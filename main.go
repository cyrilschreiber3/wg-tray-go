package main

import (
	_ "embed"
	"log/slog"
	"os"
	"os/signal"
	"strings"

	"github.com/cyrilschreiber3/wg-tray-go/models"
	"github.com/cyrilschreiber3/wg-tray-go/ui"
	"github.com/cyrilschreiber3/wg-tray-go/wgutils"
	"github.com/getlantern/systray"
)

//go:embed icon.png
var iconByte []byte

func onReady(tunnels *models.TunnelItems) {
	systray.SetIcon(iconByte)

	trayManager := ui.NewTrayManager(tunnels, handleTunnelToggle, handleUpAll, handleDownAll)
	trayManager.CreateTunnelItems()
	trayManager.CreateControlItems()

	for {
		select {
		case <-trayManager.RefreshClicked():
			slog.Info("Refresh clicked")
			trayManager.RefreshTunnelItems()

		case <-trayManager.QuitClicked():
			systray.Quit()
			return
		}
	}
}

func handleTunnelToggle(name string, active bool) {
	slog.Info("Toggling tunnel", slog.String("name", name), slog.Bool("active", active))
	status := "disabled"
	if active {
		status = "enabled"
		err := wgutils.ActivateWgTunnel(name)
		if err != nil {
			slog.Error("Error activating tunnel", slog.String("name", name), slog.Any("error", err))
		}
	} else {
		err := wgutils.DeactivateWgTunnel(name)
		if err != nil {
			slog.Error("Error deactivating tunnel", slog.String("name", name), slog.Any("error", err))
		}
	}
	slog.Info("Tunnel toggled", slog.String("name", name), slog.String("status", status))
}

func handleUpAll(tunnels []string) {
	slog.Info("Activating all tunnels")
	for _, tunnel := range tunnels {
		err := wgutils.ActivateWgTunnel(tunnel)
		if err != nil {
			slog.Error("Error activating tunnel", slog.String("name", tunnel), slog.Any("error", err))
		}
	}
}

func handleDownAll(tunnels []string) {
	slog.Info("Deactivating all tunnels")
	for _, tunnel := range tunnels {
		err := wgutils.DeactivateWgTunnel(tunnel)
		if err != nil {
			slog.Error("Error deactivating tunnel", slog.String("name", tunnel), slog.Any("error", err))
		}
	}
}

func slogLevelFromEnv() slog.Level {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("LOGLEVEL"))) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slogLevelFromEnv()}))
	slog.SetDefault(logger)
	slog.Info("Starting wg-tray-go")

	tunnels, err := wgutils.GetWgAvailableTunnels()
	if err != nil {
		slog.Error("Error getting tunnels", slog.Any("error", err))
		return
	}

	systray.Run(func() { onReady(tunnels) }, nil)
}

func init() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		slog.Info("Received interrupt signal, quitting...")
		systray.Quit()
	}()
}
