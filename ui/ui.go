package ui

import (
	"log/slog"
	"math/rand"

	"github.com/cyrilschreiber3/wg-tray-go/config"
	"github.com/cyrilschreiber3/wg-tray-go/models"
	"github.com/cyrilschreiber3/wg-tray-go/wgutils"
	"github.com/getlantern/systray"
)

type TrayManager struct {
	appConfig       *config.AppConfig
	tunnels         *models.TunnelItems
	menuItems       []*systray.MenuItem
	menuItemTunnels map[*systray.MenuItem]*models.TunnelItem
	groupMenus      map[string]*systray.MenuItem
	upAllItem       *systray.MenuItem
	downAllItem     *systray.MenuItem
	refreshItem     *systray.MenuItem
	quitItem        *systray.MenuItem
	onTunnelToggle  func(string, bool)
	onUpAll         func([]string)
	onDownAll       func([]string)
}

func NewTrayManager(appConfig *config.AppConfig, tunnels *models.TunnelItems, onTunnelToggle func(string, bool), onUpAll func([]string), onDownAll func([]string)) *TrayManager {
	return &TrayManager{
		appConfig:       appConfig,
		tunnels:         tunnels,
		menuItems:       make([]*systray.MenuItem, 0),
		menuItemTunnels: make(map[*systray.MenuItem]*models.TunnelItem),
		groupMenus:      make(map[string]*systray.MenuItem, 0),
		onTunnelToggle:  onTunnelToggle,
		onUpAll:         onUpAll,
		onDownAll:       onDownAll,
	}
}

func (tm *TrayManager) CreateTunnelItems() {
	// Add ungrouped tunnels at the top level
	ungroupedTunnels := tm.appConfig.GetUngroupedTunnelNames()
	for _, tunnelName := range ungroupedTunnels {
		tunnel := tm.tunnels.GetByName(tunnelName)
		if tunnel == nil {
			continue
		}
		menuItem := systray.AddMenuItemCheckbox(tunnel.Name, "", tunnel.Active)
		tm.menuItems = append(tm.menuItems, menuItem)
		tm.menuItemTunnels[menuItem] = tunnel

		go tm.handleTunnelClick(menuItem, tunnel)
	}

	if tm.appConfig.HasGroups() {
		systray.AddSeparator()
		for _, group := range tm.appConfig.TunnelGroups {
			groupMenu := systray.AddMenuItem(group.Name, "")
			tm.groupMenus[group.Name] = groupMenu

			// Add tunnels within this group
			for _, tunnelName := range group.TunnelNames {
				tunnel := tm.tunnels.GetByName(tunnelName)
				if tunnel == nil {
					continue
				}
				menuItem := groupMenu.AddSubMenuItemCheckbox(tunnel.Name, "", tunnel.Active)
				tm.menuItems = append(tm.menuItems, menuItem)
				tm.menuItemTunnels[menuItem] = tunnel
				go tm.handleTunnelClick(menuItem, tunnel)
			}

			upGroupText := "Up all in group"
			if group.PickRandomly {
				upGroupText = "Up random in group"
			}
			systray.AddSeparator()
			upGroupItem := groupMenu.AddSubMenuItem(upGroupText, "")
			downGroupItem := groupMenu.AddSubMenuItem("Down all in group", "")

			go tm.handleGroupUp(group, upGroupItem)
			go tm.handleGroupDown(group, downGroupItem)
		}
	}
}

func (tm *TrayManager) handleTunnelClick(menuItem *systray.MenuItem, tunnel *models.TunnelItem) {
	for range menuItem.ClickedCh {
		tunnel.ToggleActive()
		tm.updateMenuItem(menuItem, tunnel.Active)

		if tm.onTunnelToggle != nil {
			tm.onTunnelToggle(tunnel.Name, tunnel.Active)
			tm.RefreshTunnelItems()
		}
	}
}

func (tm *TrayManager) handleGroupUp(group models.TunnelGroup, upGroupItem *systray.MenuItem) {
	for range upGroupItem.ClickedCh {
		if group.PickRandomly {
			tm.RefreshTunnelItems()
			activeTunnels := tm.tunnels.GetActiveTunnelNamesInGroup(group)
			if len(activeTunnels) > 0 {
				slog.Warn("At least one tunnel in the group is already active; skipping random selection", slog.String("group", group.Name))
				continue
			}

			tunnelIndex := rand.Intn(len(group.TunnelNames))
			selectedTunnel := group.TunnelNames[tunnelIndex]
			slog.Info("Randomly selected tunnel to activate", slog.String("group", group.Name), slog.String("tunnel", selectedTunnel))
			if tm.onTunnelToggle != nil {
				tm.onTunnelToggle(selectedTunnel, true)
				tm.RefreshTunnelItems()
			}

		} else {
			tm.RefreshTunnelItems()
			if tm.onUpAll != nil {
				tm.onUpAll(group.TunnelNames)
			}
			tm.RefreshTunnelItems()
		}
	}
}

func (tm *TrayManager) handleGroupDown(group models.TunnelGroup, downGroupItem *systray.MenuItem) {
	for range downGroupItem.ClickedCh {
		tm.RefreshTunnelItems()
		if tm.onDownAll != nil {
			tm.onDownAll(tm.tunnels.GetActiveTunnelNamesInGroup(group))
		}
		tm.RefreshTunnelItems()
	}
}

func (tm *TrayManager) updateMenuItem(menuItem *systray.MenuItem, active bool) {
	if active {
		menuItem.Check()
	} else {
		menuItem.Uncheck()
	}
}

func (tm *TrayManager) RefreshTunnelItems() {
	wgutils.RefreshWgTunnels(tm.tunnels)
	for _, menuItem := range tm.menuItems {
		if tunnel, exists := tm.menuItemTunnels[menuItem]; exists {
			tm.updateMenuItem(menuItem, tunnel.Active)
		}
	}
}

func (tm *TrayManager) CreateControlItems() {
	systray.AddSeparator()
	tm.upAllItem = systray.AddMenuItem("Up all interfaces", "")
	tm.downAllItem = systray.AddMenuItem("Down all interfaces", "")
	systray.AddSeparator()
	tm.refreshItem = systray.AddMenuItem("Refresh", "Refresh tunnel states")
	tm.quitItem = systray.AddMenuItem("Quit", "Quit wg-tray-go")

	go func() {
		for range tm.upAllItem.ClickedCh {
			tm.RefreshTunnelItems()
			if tm.onUpAll != nil {
				tm.onUpAll(tm.tunnels.GetInactiveTunnelNames())
			}
			tm.RefreshTunnelItems()
		}
	}()

	go func() {
		for range tm.downAllItem.ClickedCh {
			tm.RefreshTunnelItems()
			if tm.onDownAll != nil {
				tm.onDownAll(tm.tunnels.GetActiveTunnelNames())
			}
			tm.RefreshTunnelItems()
		}
	}()
}

func (tm *TrayManager) RefreshClicked() <-chan struct{} {
	return tm.refreshItem.ClickedCh
}

func (tm *TrayManager) QuitClicked() <-chan struct{} {
	return tm.quitItem.ClickedCh
}
