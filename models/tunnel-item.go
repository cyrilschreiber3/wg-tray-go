package models

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

func (t *TunnelItems) GetInactiveTunnelNames() []string {
	inactiveTunnels := []string{}
	for i := range *t {
		if !(*t)[i].Active {
			inactiveTunnels = append(inactiveTunnels, (*t)[i].Name)
		}
	}
	return inactiveTunnels
}

func (t *TunnelItems) RefreshTrayItems() {
}
