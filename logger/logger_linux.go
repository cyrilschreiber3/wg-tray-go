package logger

import (
	"os"
	"path/filepath"
)

func getLogPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "/tmp/wg-tray-go.log"
	}
	return filepath.Join(home, ".local", "state", "wg-tray-go.log")
}
