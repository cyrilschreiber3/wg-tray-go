package wgutils

import (
	"fmt"
	"strings"
)

func parseErrorFromWg(wgOutput string) error {
	if strings.Contains(wgOutput, "Unable to access interface: No such file or directory") {
		return fmt.Errorf("interface not found")
	}

	if strings.Contains(wgOutput, "Permission denied") {
		return fmt.Errorf("permission denied")
	}

	return fmt.Errorf("unknown error: %s", wgOutput)
}

func parseErrorFromWgQuick(wgQuickOutput string) error {
	if strings.Contains(wgQuickOutput, "is not a WireGueard interface") {
		return fmt.Errorf("tunnel is not up")
	}

	if strings.Contains(wgQuickOutput, "does not exist") {
		return fmt.Errorf("tunnel configuration not found")
	}

	if strings.Contains(wgQuickOutput, "already exists as") {
		return fmt.Errorf("tunnel already active")
	}

	if strings.Contains(wgQuickOutput, "Permission denied") {
		return fmt.Errorf("permission denied")
	}

	return fmt.Errorf("unknown error: %s", wgQuickOutput)
}
