package tui

import (
	"github.com/TilliboyF/tuido/common"
	"github.com/TilliboyF/tuido/db"
	"github.com/TilliboyF/tuido/types"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))
	columns = []table.Column{
		{Title: "ID", Width: 2},
		{Title: "Name", Width: 20},
		{Title: "Status", Width: 6},
		{Title: "CreatedAt", Width: 15},
	}
)

type Model struct {
	quitting bool
	store    *db.SqliteTodoStore
	table    table.Model
}

func NewModel(store *db.SqliteTodoStore) Model {
	todos, _ := store.GetAll()
	var rows []table.Row
	for _, todo := range todos {
		rows = append(rows, common.StringArray(todo))
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithWidth(50),
		table.WithHeight(len(todos)),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	t.SetStyles(s)

	return Model{
		store: store,
		table: t,
	}

}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "e":
			f := NewForm(types.Todo{Name: "test"}, &m)
			return f.Update(nil)

		case "q":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			return m, tea.Batch(tea.Printf("Todo %s!", m.table.SelectedRow()[1]))
		}

	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}
	return baseStyle.Render(m.table.View()) + "\n  " + m.table.HelpView() + "\n"
}
