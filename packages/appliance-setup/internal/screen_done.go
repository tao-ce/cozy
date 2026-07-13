package internal

import (
	"fmt"

	"charm.land/huh/v2"
)

func (m *Model) DoneScreen() *Screen {
	configPath := m.ConfigPath
	if configPath == "" {
		configPath = m.I8nGet("done.standard_output")
	}
	return &Screen{
		Help: "done.help",
		Form: huh.NewForm(
			huh.NewGroup(
				huh.NewNote().
					Description(fmt.Sprintf(m.I8nGet("done.description"), configPath)),
				huh.NewConfirm().
					Affirmative(m.I8nGet("done.confirm.affirmative")).
					Negative(m.I8nGet("done.confirm.negative")).
					Title(m.I8nGet("done.confirm.title")).
					Value(&m.Restart),
			),
		),
	}
}
