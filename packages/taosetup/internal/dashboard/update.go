package dashboard

import (
	"time"

	tea "charm.land/bubbletea/v2"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Window.Width = msg.Width
		m.Window.Height = msg.Height
		return m, nil
	case HealthcheckMsg:
		cmds = append(cmds, tea.Every(2*time.Second, m.updateHealthcheck))

	case tea.KeyPressMsg:
		cmds = append(cmds, tea.ClearScreen)

		switch msg.String() {
		case "ctrl+c":
			return m, m.quit()
		}
	}

	_, cmd = m.MenuWidget.Update(msg)
	cmds = append(cmds, cmd)

	// m.HealthcheckWidget.table, cmd = m.HealthcheckWidget.table.Update(msg)
	// cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
