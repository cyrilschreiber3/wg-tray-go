package main

import (
	_ "embed"
	"fmt"

	"github.com/getlantern/systray"
)

//go:embed icon.png
var iconByte []byte

var strayUpAllItem *systray.MenuItem
var strayRefreshItem *systray.MenuItem
var strayQuitItem *systray.MenuItem

type TunnelItem struct {
	Name         string
	Active       bool
	TrayMenuItem *systray.MenuItem
}

var tunnelItems = []TunnelItem{
	{"Tunnel 1", true, nil},
	{"Tunnel 2", false, nil},
}

func toggleTunnelItem(tunnelItemId int) {
	tunnelItem := tunnelItems[tunnelItemId]

	if tunnelItem.Active {
		fmt.Printf("%s was on, disabling...", tunnelItem.Name)
		tunnelItems[tunnelItemId].Active = false
		tunnelItem.TrayMenuItem.Uncheck()
	} else {
		fmt.Printf("%s was off, activating...", tunnelItem.Name)
		tunnelItems[tunnelItemId].Active = true
		tunnelItem.TrayMenuItem.Check()
	}
}

func listTunnelItems() {

	// getTunnelData()

	for id, tunnel := range tunnelItems {
		strayTunnelItem := systray.AddMenuItemCheckbox(tunnel.Name, "", tunnel.Active)
		tunnelItems[id].TrayMenuItem = strayTunnelItem

		go func(item *systray.MenuItem) {
			for {
				<-item.ClickedCh
				fmt.Printf("%s clicked", tunnel.Name)
				toggleTunnelItem(id)
			}
		}(strayTunnelItem)
	}
}

func refreshTunnelItems() {

	// updateTunnelData()

	for _, item := range tunnelItems {
		if item.Active {
			item.TrayMenuItem.Check()
		} else {
			item.TrayMenuItem.Uncheck()
		}
	}

}

func listAllItems() {
	listTunnelItems()

	systray.AddSeparator()
	strayUpAllItem = systray.AddMenuItem("Up all interfaces", "")
	systray.AddSeparator()
	strayRefreshItem = systray.AddMenuItem("Refresh", "Refresh the list of tunnels")
	strayQuitItem = systray.AddMenuItem("Quit", "Quit wg-menu-bar")
}

func onReady() {
	systray.SetIcon(iconByte)

	listAllItems()

	for {
		select {
		case <-strayUpAllItem.ClickedCh:
			fmt.Println("Up all clicked")
			for id, tunnel := range tunnelItems {
				if !tunnel.Active {
					toggleTunnelItem(id)
				}
			}

		case <-strayRefreshItem.ClickedCh:
			fmt.Println("Refresh clicked")
			refreshTunnelItems()

		case <-strayQuitItem.ClickedCh:
			systray.Quit()
			return
		}
	}
}

func main() {
	systray.Run(onReady, nil)
}
