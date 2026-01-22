package wgutils

import (
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/cyrilschreiber3/wg-tray-go/models"
)

var wgConfigPath = "/etc/wireguard/"

func GetWgAvailableTunnels() (*models.TunnelItems, error) {
	tunnelsConfs, err := os.ReadDir(wgConfigPath)
	if err != nil {
		slog.Error("Error reading wg config directory", slog.Any("error", err))
		return nil, err
	}

	tunnels := &models.TunnelItems{}

	for _, file := range tunnelsConfs {
		if file.IsDir() {
			continue
		}
		name := strings.TrimSuffix(file.Name(), ".conf")
		active := IsWgTunnelActive(name)
		*tunnels = append(*tunnels, models.TunnelItem{Name: name, Active: active})
	}

	return tunnels, nil
}

func IsWgTunnelActive(tunnelName string) bool {
	interfaceName, err := getInterfaceName(tunnelName)
	if err != nil {
		return false
	}

	output, err := exec.Command("wg", "show", interfaceName).CombinedOutput()
	slog.Debug("wg show output", slog.String("output", string(output)))
	if err != nil {
		errorString := parseErrorFromWg(string(output))
		slog.Error("Error checking tunnel status", slog.String("interface", interfaceName), slog.Any("error", errorString))
		return false
	}

	return len(output) > 0
}

func ActivateWgTunnel(tunnelName string) error {
	output, err := exec.Command("wg-quick", "up", tunnelName).CombinedOutput()
	slog.Debug("wg-quick up output", slog.String("output", string(output)))
	if err != nil {
		errorString := parseErrorFromWgQuick(string(output))
		return errorString
	}
	return nil
}

func DeactivateWgTunnel(tunnelName string) error {
	output, err := exec.Command("wg-quick", "down", tunnelName).CombinedOutput()
	slog.Debug("wg-quick down output", slog.String("output", string(output)))
	if err != nil {
		errorString := parseErrorFromWgQuick(string(output))
		return errorString
	}
	return nil
}

func RefreshWgTunnels(tunnels *models.TunnelItems) {
	for index, tunnel := range *tunnels {
		(*tunnels)[index].Active = IsWgTunnelActive(tunnel.Name)
	}
}
