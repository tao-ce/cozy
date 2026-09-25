package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/tao-ce/cozy/taosetup/internal"
	"github.com/tao-ce/cozy/taosetup/internal/dashboard"
	"github.com/tao-ce/cozy/taosetup/internal/install"
)

//go:embed default.json
var defaultConfig string

func main() {

	subcommand := "dashboard"

	if len(os.Args) > 1 {
		subcommand = os.Args[1]
	} else {
		fmt.Fprintln(os.Stderr, "missing subcommand, running dashboard")
	}

	installCmd := flag.NewFlagSet("install", flag.ExitOnError)
	dashboardCmd := flag.NewFlagSet("dashboard", flag.ExitOnError)

	installConfigPath := installCmd.String("config-path", "", "path to the config file")
	dashboardConfigPath := dashboardCmd.String("config-path", "", "path to the config file")

	switch subcommand {
	case "install":
		installCmd.Parse(os.Args[2:])
		modeInstall(*installConfigPath)
	case "dashboard":
		if len(os.Args) > 1 {
			dashboardCmd.Parse(os.Args[2:])
		}
		modeDashboard(*dashboardConfigPath)
	default:
		fmt.Fprintln(os.Stderr, "invalid subcommand:", subcommand)
		os.Exit(1)
	}
}

func modeDashboard(configPath string) {
	var setup internal.ApplianceSetup
	setup.Import(configPath)

	model := dashboard.Model{
		Setup:      setup,
		ConfigPath: configPath,
	}
	p := tea.NewProgram(&model)
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "could not start program:", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func modeInstall(configPath string) {
	var setup internal.ApplianceSetup

	if err := json.Unmarshal([]byte(defaultConfig), &setup); err != nil {
		fmt.Fprintln(os.Stderr, "could not unmarshal default config:", err)
		os.Exit(1)
	}

	model := install.Model{
		Setup:      &setup,
		ConfigPath: configPath,
	}

	p := tea.NewProgram(&model)
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "could not start program:", err)
		os.Exit(1)
	}
	model.Setup.Export(configPath)
	os.Exit(0)
}
