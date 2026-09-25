package menu

import (
	"fmt"
	"io"
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type item struct {
	title    string
	index    int
	subitems []list.Item
}

var menuItems = []list.Item{
	item{title: "Healthcheck"},
	item{title: "TAO", subitems: []list.Item{
		item{title: "Admins", subitems: []list.Item{
			item{title: "Add Admin"},
			item{title: "Remove Admin"},
		}},
		item{title: "Users"},
	}},
	item{title: "System", subitems: []list.Item{
		item{title: "Shutdown"},
		item{title: "Reboot"},
	}},
	item{title: "Exit"},
}

type styles struct {
	item     lipgloss.Style
	path     lipgloss.Style
	selected lipgloss.Style
}
type itemDelegate struct {
	styles *styles
}

func (i item) Title() string        { return i.title }
func (i item) FilterValue() string  { return i.title }
func (i item) Index() int           { return i.index }
func (i itemDelegate) Height() int  { return 1 }
func (i itemDelegate) Spacing() int { return 0 }

func (i itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	fn := d.styles.item.Render
	if index == m.Index() {
		fn = d.styles.selected.Render
	}

	fmt.Fprint(w, fn(i.title))
}

type MenuWidget struct {
	menu         list.Model
	previousMenu *[]list.Item
	path         []string
	selected     int
	active       int
}

func New() *MenuWidget {

	menu := list.New(
		menuItems,
		itemDelegate{
			styles: &styles{
				item:     lipgloss.NewStyle().PaddingLeft(4),
				path:     lipgloss.NewStyle().PaddingLeft(4),
				selected: lipgloss.NewStyle().PaddingLeft(4).Foreground(lipgloss.Color("226")).Bold(true),
			},
		},
		40,
		10,
	)

	menu.Title = "Menu"

	menu.Styles.Title = lipgloss.NewStyle().Background(lipgloss.Color("239")).Bold(true)
	menu.SetShowFilter(false)
	menu.SetShowHelp(false)
	menu.SetShowStatusBar(false)
	// menu.SetShowTitle(false)
	menu.SetShowPagination(false)

	return &MenuWidget{
		selected: 0,
		active:   0,
		menu:     menu,
		path:     []string{"Menu"},
	}
}

func (w *MenuWidget) Update(msg tea.Msg) (*MenuWidget, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up":
			w.menu.CursorUp()
			w.selected = w.menu.Index()
		case "down":
			w.menu.CursorDown()
			w.selected = w.menu.Index()
		case "enter":
			if w.menu.Items()[w.menu.Index()].(item).subitems != nil {
				w.path = append(w.path, w.menu.Items()[w.menu.Index()].(item).title)
				previousMenu := w.menu.Items()
				w.previousMenu = &previousMenu
				w.menu.SetItems(w.menu.Items()[w.menu.Index()].(item).subitems)
				w.menu.InsertItem(99, item{title: "Back"})
				w.menu.Select(0)
				// w.active = 0
				w.menu.Title = strings.Join(w.path, " > ")
			} else {
				w.active = w.menu.Index()
				if w.menu.SelectedItem().FilterValue() == "Back" {
					w.path = w.path[:len(w.path)-1]
					w.menu.SetItems(*w.previousMenu)
					w.menu.Select(0)
					w.menu.Title = strings.Join(w.path, " > ")
					w.active = w.selected
					w.selected = 0
				}
			}
		}
	}

	_, cmd := w.menu.Update(msg)
	return w, cmd
}

func (w *MenuWidget) View() string {
	return w.menu.View()
}

func (w *MenuWidget) ActiveItem() string {
	return w.menu.Items()[w.active].(item).title
}

func (w *MenuWidget) SelectedItem() string {
	return w.menu.SelectedItem().FilterValue()
}
