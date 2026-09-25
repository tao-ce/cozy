package dashboard

import (
	_ "embed"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	"github.com/tao-ce/cozy/taosetup/internal/i18n"
)

// //go:embed taobanner.txt
// var taobanner string

// //go:embed taoname.txt
// var taoname string

//go:embed taologo.ansi
var taologo string

type Screen struct {
	Help string
	Form *huh.Form
}

func (m *Model) I8nGet(key string) string {
	return i18n.Get("en", key)
}

func (m *Model) selectedTab() string {
	switch m.MenuWidget.ActiveItem() {
	case "Healthcheck":
		return m.renderHealthcheckTable()
	// case "System":
	// return getUptime().Round(time.Second).String()
	// case "Shutdown":
	// 	return "shutdown"
	// case "Reboot":
	// 	return "reboot"
	case "Exit":
		return "exit"
	}
	return ""
}

func (m *Model) View() tea.View {
	// screen := m.DashboardScreen()

	leftPaneWidth := int(float64(m.Window.Width) * 0.37)
	rightPaneWidth := int(float64(m.Window.Width) * 0.63)

	leftPaneStyle := lipgloss.NewStyle().
		Width(leftPaneWidth).
		Height(m.Window.Height - 1).
		Align(lipgloss.Center).
		Border(lipgloss.HiddenBorder()).
		Padding(1)
	rightPaneStyle := lipgloss.NewStyle().
		Width(rightPaneWidth).
		Height(m.Window.Height - 1).
		Border(lipgloss.NormalBorder()).
		Padding(1).
		Align(lipgloss.Center)
	// helpPaneStyle := lipgloss.NewStyle().
	// 	Width(rightPaneWidth).
	// 	Height(1).
	// 	Align(lipgloss.Left).
	// 	Padding(0, 0, 0, 1)

	v := tea.NewView("")

	v.Content = lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			leftPaneStyle.Render(
				lipgloss.Place(42, 0, lipgloss.Left, lipgloss.Top, taologo),
				"\n",
				lipgloss.Place(42, 0, lipgloss.Left, lipgloss.Top, m.MenuWidget.View()),
			),
			rightPaneStyle.Render(
				m.selectedTab(),
				// fmt.Sprintf("Uptime: %s", getUptime().Round(time.Second).String()),
			),
		),
		// helpPaneStyle.Render(screen.Form.Help().ShortHelpView(screen.Form.KeyBinds())),
	)
	v.AltScreen = true
	return v
}
