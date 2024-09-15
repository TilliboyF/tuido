package tui

import (
	"strconv"

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
			BorderForeground(lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"})
	columns = []table.Column{
		{Title: "ID", Width: 2},
		{Title: "Name", Width: 15},
		{Title: "Status", Width: 12},
		{Title: "CreatedAt", Width: 25},
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
		// table.WithWidth(50),
		table.WithHeight(len(todos)+5),
	)

	s := table.DefaultStyles()

	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}).
		BorderBottom(true).
		Bold(false).
		Foreground(lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"})
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

	case types.Todo:
		m.table, cmd = m.HandleTodoReturn(msg)
		return m, cmd

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "e":
			currentRow := m.table.SelectedRow()
			todoID, err := strconv.ParseInt(currentRow[0], 10, 64)
			if err != nil {
				panic("That should not happen")
			}

			todo, err := m.store.GetById(todoID)
			if err != nil {
				panic("That should not happen")
			}

			f := NewForm(todo, &m)
			return f.Update(nil)

		case "n":
			f := NewEmptyForm(&m)
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

func (m Model) HandleTodoReturn(todo types.Todo) (table.Model, tea.Cmd) {
	if todo.ID == -1 { // new todo
		m.store.Add(&todo)

		todo, _ = m.store.GetById(todo.ID)

		currentRows := m.table.Rows()
		currentRows = append(currentRows, common.StringArray(todo))
		m.table.SetRows(currentRows)
		m.table.UpdateViewport()
		return m.table.Update(nil)
	} else { // existing one
		err := m.store.Update(&todo)
		if err != nil {
			//tbd
		}
		todos, _ := m.store.GetAll()
		var rows []table.Row
		for _, todo := range todos {
			rows = append(rows, common.StringArray(todo))
		}
		m.table.SetRows(rows)
		m.table.UpdateViewport()
		return m.table.Update(nil)
	}
}
