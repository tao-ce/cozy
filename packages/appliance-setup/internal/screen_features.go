package internal

import (
	"fmt"

	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

func (m *Model) FeaturesScreen() *Screen {
	warning := "\n"
	if !m.Is(IsWiFiHotspotAvailable) {
		warning += lipgloss.NewStyle().Foreground(lipgloss.Yellow).Bold(true).Render(
			fmt.Sprintf(
				m.I8nGet("features.wifihotspot.warning"),
				m.I8nGet(fmt.Sprintf("medium.%s", m.Setup.Medium)),
			))
	}
	return &Screen{
		Help: "features.help",
		Form: huh.NewForm(
			huh.NewGroup(
				huh.NewNote().Description(warning),

				huh.NewConfirm().
					Title(fmt.Sprintf("%-25s", m.I8nGet("features.wifihotspot"))).
					Value(&m.Setup.ApplianceFeatures.StartWiFiHotspot).
					Inline(true).
					Affirmative(m.I8nGet("features.on")).
					Negative(m.I8nGet("features.off")),

				huh.NewConfirm().
					Title(fmt.Sprintf("%-25s", m.I8nGet("features.multicastdns"))).
					Value(&m.Setup.ApplianceFeatures.StartMulticastDNS).
					Inline(true).
					Affirmative(m.I8nGet("features.on")).
					Negative(m.I8nGet("features.off")),
			),
		),
	}

}
