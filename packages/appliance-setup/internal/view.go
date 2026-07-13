package internal

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/glamour/v2"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

type Screen struct {
	Help string
	Form *huh.Form
}

const (
	minTerminalWidth  = 80
	minTerminalHeight = 20
)

func (m *Model) View() tea.View {
	screen := m.GetScreen()

	leftPaneWidth := int(float64(m.WindowWidth) * 0.37)
	rightPaneWidth := int(float64(m.WindowWidth) * 0.63)

	leftPaneStyle := lipgloss.NewStyle().
		Width(leftPaneWidth).
		Height(m.WindowHeight - 1).
		Border(lipgloss.NormalBorder()).
		Padding(1)
	rightPaneStyle := lipgloss.NewStyle().
		Width(rightPaneWidth).
		Height(m.WindowHeight - 1).
		Border(lipgloss.NormalBorder()).
		Padding(1).
		Align(lipgloss.Center)
	helpPaneStyle := lipgloss.NewStyle().
		Width(rightPaneWidth).
		Height(1).
		Align(lipgloss.Left).
		Padding(0, 0, 0, 1)

	helpRenderer, _ := glamour.NewTermRenderer(
		glamour.WithWordWrap(leftPaneWidth-4),
		glamour.WithStandardStyle("dark"),
	)

	if screen == nil {
		return tea.NewView("")
	}

	help, _ := helpRenderer.Render(m.I8nGet(screen.Help))

	formView := screen.Form.
		WithLayout(huh.LayoutStack).
		WithShowHelp(false).
		View()

	v := tea.NewView(
		lipgloss.JoinVertical(
			lipgloss.Top,
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				leftPaneStyle.Render(
					help,
				),
				rightPaneStyle.Render(formView),
			),
			helpPaneStyle.Render(screen.Form.Help().ShortHelpView(screen.Form.KeyBinds())),
		),
	)
	v.AltScreen = true
	return v
}
