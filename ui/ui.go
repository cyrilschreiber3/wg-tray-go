package ui

import (
	"github.com/cyrilschreiber3/wg-tray-go/models"
	"github.com/cyrilschreiber3/wg-tray-go/wgutils"
	"github.com/getlantern/systray"
)

type TrayManager struct {
	tunnels        *models.TunnelItems
	menuItems      []*systray.MenuItem
	upAllItem      *systray.MenuItem
	downAllItem    *systray.MenuItem
	refreshItem    *systray.MenuItem
	quitItem       *systray.MenuItem
	onTunnelToggle func(string, bool)
	onUpAll        func([]string)
	onDownAll      func([]string)
}

func NewTrayManager(tunnels *models.TunnelItems, onTunnelToggle func(string, bool), onUpAll func([]string), onDownAll func([]string)) *TrayManager {
	return &TrayManager{
		tunnels:        tunnels,
		menuItems:      make([]*systray.MenuItem, 0),
		onTunnelToggle: onTunnelToggle,
		onUpAll:        onUpAll,
		onDownAll:      onDownAll,
	}
}

func (tm *TrayManager) CreateTunnelItems() {
	for i := range *tm.tunnels {
		tunnel := &(*tm.tunnels)[i]
		menuItem := systray.AddMenuItemCheckbox(tunnel.Name, "", tunnel.Active)
		tm.menuItems = append(tm.menuItems, menuItem)

		go tm.handleTunnelClick(menuItem, tunnel)
	}
}

func (tm *TrayManager) handleTunnelClick(menuItem *systray.MenuItem, tunnel *models.TunnelItem) {
	for range menuItem.ClickedCh {
		tunnel.ToggleActive()
		tm.updateMenuItem(menuItem, tunnel.Active)

		if tm.onTunnelToggle != nil {
			tm.onTunnelToggle(tunnel.Name, tunnel.Active)
		}
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
	for i, menuItem := range tm.menuItems {
		if i < len(*tm.tunnels) {
			tm.updateMenuItem(menuItem, (*tm.tunnels)[i].Active)
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
