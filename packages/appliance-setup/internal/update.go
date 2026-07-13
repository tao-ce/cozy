package internal

import (
	"time"

	"charm.land/bubbles/v2/timer"
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.WindowWidth = msg.Width
		m.WindowHeight = msg.Height
		return m, nil
	case timer.TickMsg:
		var cmd tea.Cmd
		m.Timer, cmd = m.Timer.Update(msg)
		cmds = append(cmds, cmd)
	case timer.TimeoutMsg:
		if m.SetupStep == SetupStepWelcome {
			return m, m.quit()
		}
	case tea.KeyPressMsg:
		cmds = append(cmds, tea.ClearScreen)
		if m.SetupStep == SetupStepWelcome && m.Timer.Running() {
			m.Timer.Timeout = time.Duration(0)
		}
		switch msg.String() {
		case "esc":
			m.restart()
			return m, m.GetScreen().Form.Init()
		case "ctrl+c":
			return m, m.quit()
		}
	}

	screen := m.GetScreen()
	if screen != nil {
		form, cmd := screen.Form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			screen.Form = f
			cmds = append(cmds, cmd)
		}
	}

	if screen.Form.State == huh.StateCompleted {
		if m.UseDefaultSettings {
			return m, m.quit()
		}

		if m.Restart {
			m.restart()
			return m, m.GetScreen().Form.Init()
		}

		m.SetupStep++

		if m.SetupStep == SetupStepWiFiHotspotSetup && !m.Setup.ApplianceFeatures.StartWiFiHotspot {
			m.SetupStep++
		}

		if m.SetupStep > SetupStepDone {
			return m, m.quit()
		}

		screen = m.GetScreen()
		if screen != nil {
			cmds = append(cmds, screen.Form.Init())
		}
	}
	return m, tea.Batch(cmds...)
}
