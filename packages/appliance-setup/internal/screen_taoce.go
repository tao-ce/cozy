package internal

import (
	"fmt"

	"charm.land/huh/v2"
)

func (m *Model) TAOCEConfigScreen() *Screen {
	return &Screen{
		Help: "taoce.help",
		Form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title(m.I8nGet("taoce.fqdn")).
					CharLimit(80).
					Value(&m.Setup.TAOCEConfig.FQDN).Validate(m.Validate("taoce.fqdn")).
					WithWidth(32),
				huh.NewSelect[string]().
					Title(m.I8nGet("taoce.flavor")).
					Value(&m.Setup.TAOCEConfig.Flavor).
					Options(
						huh.NewOption(m.I8nGet("taoce.flavor.full"), "full"),
						huh.NewOption(m.I8nGet("taoce.flavor.essential"), "essential"),
						huh.NewOption(m.I8nGet("taoce.flavor.lite"), "lite"),
						huh.NewOption(m.I8nGet("taoce.flavor.minimal"), "minimal"),
					).
					WithWidth(32),
				huh.NewNote().DescriptionFunc(func() string {
					return m.I8nGet(fmt.Sprintf("taoce.flavor.%s.description", m.Setup.TAOCEConfig.Flavor))
				}, &m.Setup.TAOCEConfig.Flavor),
			),
		),
	}
}
