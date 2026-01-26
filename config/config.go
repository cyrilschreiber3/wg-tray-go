package config

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cyrilschreiber3/wg-tray-go/models"
	"github.com/cyrilschreiber3/wg-tray-go/wgutils"
)

var wgConfigPath = "/etc/wireguard/"

type AppConfig struct {
	TunnelNames  []string             `json:"tunnel_names"`
	TunnelGroups []models.TunnelGroup `json:"tunnel_groups"`
}

func NewAppConfig() *AppConfig {
	return &AppConfig{
		TunnelNames:  []string{},
		TunnelGroups: []models.TunnelGroup{},
	}
}

func getConfigPath() string {
	if p := os.Getenv("WG_TRAY_CONFIG"); p != "" {
		return p
	}
	dir := os.Getenv("XDG_CONFIG_HOME")
	if dir != "" {
		return filepath.Join(dir, "wg-tray-go", "config.json")
	}
	home, err := os.UserHomeDir()
	if err == nil && home != "" {
		return filepath.Join(home, ".config", "wg-tray-go", "config.json")
	}
	return filepath.Join(".config", "wg-tray-go", "config.json")
}

func loadConfigFromFile(path string) (*AppConfig, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			appConfig, terr := getTunnelsFromWgConfig()
			if terr != nil {
				return nil, terr
			}
			return appConfig, nil
		}
		return nil, err
	}

	appConfig := NewAppConfig()

	err = json.Unmarshal(file, &appConfig)
	if err != nil {
		return nil, err
	}

	return appConfig, nil
}

func getTunnelsFromWgConfig() (*AppConfig, error) {
	output, err := exec.Command("sudo", "ls", wgConfigPath).CombinedOutput()
	if err != nil {
		return nil, err
	}

	appConfig := NewAppConfig()

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, file := range files {
		if strings.HasSuffix(file, ".conf") {
			name := strings.TrimSuffix(file, ".conf")
			appConfig.TunnelNames = append(appConfig.TunnelNames, name)
		}
	}

	return appConfig, nil
}

func LoadAppConfig() (*AppConfig, error) {
	configPath := getConfigPath()

	appConfig, err := loadConfigFromFile(configPath)
	if err != nil {
		return nil, err
	}
	return appConfig, nil
}

func (c *AppConfig) HasGroups() bool {
	return len(c.TunnelGroups) > 0
}

func (c *AppConfig) GetAllTunnelNames() []string {
	seen := make(map[string]struct{}, len(c.TunnelNames))
	names := make([]string, 0, len(c.TunnelNames))

	for _, n := range c.TunnelNames {
		if _, ok := seen[n]; !ok {
			seen[n] = struct{}{}
			names = append(names, n)
		}
	}

	for _, g := range c.TunnelGroups {
		for _, n := range g.TunnelNames {
			if _, ok := seen[n]; !ok {
				seen[n] = struct{}{}
				names = append(names, n)
			}
		}
	}

	return names
}

func (c *AppConfig) GetUngroupedTunnelNames() []string {
	return c.TunnelNames
}

func (c *AppConfig) ToTunnelItems() (*models.TunnelItems, error) {
	allTunnelNames := c.GetAllTunnelNames()

	tunnelItems := &models.TunnelItems{}

	for _, name := range allTunnelNames {
		active := wgutils.IsWgTunnelActive(name)
		*tunnelItems = append(*tunnelItems, models.TunnelItem{Name: name, Active: active})
	}

	return tunnelItems, nil
}
