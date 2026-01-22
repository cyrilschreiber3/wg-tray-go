package wgutils

import (
	"log/slog"
	"os"
	"strings"
)

var wgRuntimeDir = "/var/run/wireguard/"

func getInterfaceName(tunnelName string) (string, error) {
	interfaceName, err := os.ReadFile(wgRuntimeDir + tunnelName + ".name")
	if err != nil {
		if os.IsNotExist(err) {
			return "", err
		}
		slog.Error("Error reading interface name file", slog.Any("error", err))
		return "", err
	}

	interfaceNameStr := strings.TrimSpace(string(interfaceName))
	return interfaceNameStr, nil
}
