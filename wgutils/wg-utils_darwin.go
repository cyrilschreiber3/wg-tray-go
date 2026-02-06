package wgutils

import (
	"os/exec"
	"strings"
)

var wgRuntimeDir = "/var/run/wireguard/"

func getInterfaceName(tunnelName string) (string, error) {
	output, err := exec.Command("sudo", "cat", wgRuntimeDir+tunnelName+".name").Output()
	if err != nil {
		return "", err
	}

	interfaceNameStr := strings.TrimSpace(string(output))
	return interfaceNameStr, nil
}
