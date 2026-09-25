package dashboard

import (
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/tao-ce/cozy/taosetup/internal"
	"github.com/tao-ce/cozy/taosetup/internal/dashboard/menu"
	"github.com/tao-ce/cozy/taosetup/internal/services/tao_ce"
)

type Model struct {
	Window            Window
	HealthcheckWidget HealthcheckWidget
	MenuWidget        *menu.MenuWidget
	ConfigPath        string
	Setup             internal.ApplianceSetup
	TAOCEService      *tao_ce.TAOCEService
}

type Window struct {
	Width  int
	Height int
}

type HealthcheckWidget struct {
	table  [][]string
	report tao_ce.HealthcheckReport
}

func (m *Model) Init() tea.Cmd {
	m.TAOCEService = tao_ce.New(fmt.Sprintf("https://%s", m.Setup.TAOCEConfig.FQDN))
	m.MenuWidget = menu.New()
	m.HealthcheckWidget.table = [][]string{}
	return tea.Batch(
		tea.Every(1*time.Second, m.updateHealthcheck),
	)
}

func (m *Model) quit() tea.Cmd {
	return tea.Batch(tea.Quit)
}
