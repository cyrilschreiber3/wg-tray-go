package models

type TunnelGroup struct {
	Name         string   `json:"name"`
	PickRandomly bool     `json:"pick_randomly,omitempty"`
	TunnelNames  []string `json:"tunnel_names"`
}

type TunnelItem struct {
	Name   string
	Active bool
}

type TunnelItems []TunnelItem

func (t *TunnelItem) ToggleActive() {
	t.Active = !t.Active
}

func (t *TunnelItem) Activate() {
	t.Active = true
}

func (t *TunnelItem) Deactivate() {
	t.Active = false
}

func (t *TunnelItems) ActivateAll() {
	for i := range *t {
		(*t)[i].Active = true
	}
}

func (t *TunnelItems) DeactivateAll() {
	for i := range *t {
		(*t)[i].Active = false
	}
}

func (t *TunnelItems) GetByName(name string) *TunnelItem {
	for i := range *t {
		if (*t)[i].Name == name {
			return &(*t)[i]
		}
	}
	return nil
}

func (t *TunnelItems) GetActiveTunnelNames() []string {
	activeTunnels := []string{}
	for i := range *t {
		if (*t)[i].Active {
			activeTunnels = append(activeTunnels, (*t)[i].Name)
		}
	}
	return activeTunnels
}

func (t *TunnelItems) GetActiveTunnelNamesInGroup(group TunnelGroup) []string {
	activeTunnels := []string{}
	for _, tunnelName := range group.TunnelNames {
		tunnel := t.GetByName(tunnelName)
		if tunnel != nil && tunnel.Active {
			activeTunnels = append(activeTunnels, tunnel.Name)
		}
	}
	return activeTunnels
}

func (t *TunnelItems) GetInactiveTunnelNames() []string {
	inactiveTunnels := []string{}
	for i := range *t {
		if !(*t)[i].Active {
			inactiveTunnels = append(inactiveTunnels, (*t)[i].Name)
		}
	}
	return inactiveTunnels
}

func (t *TunnelItems) GetInactiveTunnelNamesInGroup(group TunnelGroup) []string {
	inactiveTunnels := []string{}
	for _, tunnelName := range group.TunnelNames {
		tunnel := t.GetByName(tunnelName)
		if tunnel != nil && !tunnel.Active {
			inactiveTunnels = append(inactiveTunnels, tunnel.Name)
		}
	}
	return inactiveTunnels
}
