package dashboard

import (
	"slices"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
)

type HealthcheckMsg time.Time

func (m *Model) renderHealthcheckTable() string {
	return table.New().
		Rows(m.HealthcheckWidget.table...).
		Headers("Apps", "Status").
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("167"))
			case col == 1 && m.HealthcheckWidget.table[row][1] == "ok":
				return lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Transform(func(s string) string { return " ✅ OK  " })
			case col == 1 && m.HealthcheckWidget.table[row][1] != "ok":
				return lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Transform(func(s string) string { return " ❗ Fail" })
			default:
				return lipgloss.NewStyle()
			}
		}).
		Border(lipgloss.HiddenBorder()).
		BorderColumn(false).
		Width(30).
		Render()
}

func (m *Model) updateHealthcheck(t time.Time) tea.Msg {

	rows := [][]string{}

	report, err := m.TAOCEService.Healthcheck()
	if err != nil {
		return nil
	}

	m.HealthcheckWidget.report = report

	for app, status := range m.HealthcheckWidget.report.AppsReports {
		if status > 0 {
			rows = append(rows, []string{app, "error"})
		} else {
			rows = append(rows, []string{app, "ok"})
		}
	}

	slices.SortFunc(rows, func(i, j []string) int {
		return strings.Compare(i[0], j[0])
	})

	m.HealthcheckWidget.table = rows

	return HealthcheckMsg(time.Now())

}
