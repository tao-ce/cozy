package dashboard

import (
	"charm.land/huh/v2"
)

func (m *Model) DashboardScreen() *Screen {
	return &Screen{
		Help: "dashboard.help",

		Form: huh.NewForm(
			huh.NewGroup(
				huh.NewNote().Title("dashboard"),
			)),
	}
}
