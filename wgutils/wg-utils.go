package wgutils

import (
	"log/slog"
	"os/exec"

	"github.com/cyrilschreiber3/wg-tray-go/models"
)

func IsWgTunnelActive(tunnelName string) bool {
	interfaceName, err := getInterfaceName(tunnelName)
	if err != nil {
		return false
	}

	output, err := exec.Command("sudo", "wg", "show", interfaceName).CombinedOutput()
	slog.Debug("wg show output", slog.String("output", string(output)))
	if err != nil {
		errorString := parseErrorFromWg(string(output))
		slog.Error("Error checking tunnel status", slog.String("interface", interfaceName), slog.Any("error", errorString))
		return false
	}

	return len(output) > 0
}

func ActivateWgTunnel(tunnelName string) error {
	output, err := exec.Command("sudo", "wg-quick", "up", tunnelName).CombinedOutput()
	slog.Debug("wg-quick up output", slog.String("output", string(output)))
	if err != nil {
		errorString := parseErrorFromWgQuick(string(output))
		return errorString
	}
	return nil
}

func DeactivateWgTunnel(tunnelName string) error {
	output, err := exec.Command("sudo", "wg-quick", "down", tunnelName).CombinedOutput()
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
