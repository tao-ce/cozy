package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/tao-ce/cozy/root/usr/libexec/tao-ce/cozy/appliance-setup/internal"
)

//go:embed default.json
var defaultConfig string

func main() {
	var setup internal.ApplianceSetup

	if err := json.Unmarshal([]byte(defaultConfig), &setup); err != nil {
		fmt.Fprintln(os.Stderr, "could not unmarshal default config:", err)
		os.Exit(1)
	}

	model := internal.Model{
		Setup: setup,
	}

	p := tea.NewProgram(&model)
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "could not start program:", err)
		os.Exit(1)
	}
	model.Export()

}
