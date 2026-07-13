package internal

import (
	"fmt"
	"time"

	"charm.land/bubbles/v2/progress"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

func (m *Model) WelcomeScreen() *Screen {
	return &Screen{
		Help: `# TAO CE Appliance Setup 

Follow this guide to configure your appliance.

`,
		Form: huh.NewForm(
			huh.NewGroup(
				huh.NewNote().Description(
					lipgloss.NewStyle().
						Height(2).
						Align(lipgloss.Center).
						Width(30).
						Margin(2, 0, 0).
						Padding(1, 6).
						Bold(true).
						Italic(true).
						Background(lipgloss.BrightRed).
						Render("Welcome to TAO CE\nAppliance Setup")),

				huh.NewNote().DescriptionFunc(
					func() string {
						if !m.Timer.Running() {
							return "\n"
						}

						timeoutRatio := float64(m.Timer.Timeout.Seconds()) / float64(m.WelcomeTimeout.Seconds())
						progressBar := progress.New(
							progress.WithoutPercentage(),
							progress.WithDefaultBlend(),
							// progress.WithColors(lipgloss.BrightYellow),
						)

						progressBar.SetWidth(int(float64(m.WindowWidth) * 0.5))

						return fmt.Sprintf(
							"Appliance Setup will proceed with default settings in %s",
							time.Duration(m.Timer.Timeout).Round(time.Second).String(),
						) + "\n" + lipgloss.NewStyle().Render(progressBar.ViewAs(timeoutRatio))
					},
					&m.Timer.Timeout,
				),

				huh.NewConfirm().
					Affirmative("Default settings").
					Negative("Customize settings").
					Value(&m.UseDefaultSettings).
					Title("Do you want to continue with default settings? "),
			),
		),
	}
}
