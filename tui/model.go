package tui

import (
	"strconv"

	"github.com/TilliboyF/tuido/common"
	"github.com/TilliboyF/tuido/db"
	"github.com/TilliboyF/tuido/types"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

var (
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

	styles ModelStyles

	keys Keys
	help help.Model
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
		table.WithHeight(8),
	)

	styles := NewDefaultModelStyles()

	t.SetStyles(styles.tableStyle)

	return Model{
		store:  store,
		table:  t,
		help:   help.New(),
		keys:   defaultKeys(),
		styles: styles,
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
		switch {
		case key.Matches(msg, m.keys.Edit):
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

		case key.Matches(msg, m.keys.New):
			f := NewEmptyForm(&m)
			return f.Update(nil)

		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, m.keys.Delete):
			m.table, cmd = m.HandleDeleteTask()
			return m, cmd
		case key.Matches(msg, m.keys.View):
			currentRow := m.table.SelectedRow()
			todoID, err := strconv.ParseInt(currentRow[0], 10, 64)
			if err != nil {
				panic("That should not happen")
			}

			todo, err := m.store.GetById(todoID)
			if err != nil {
				panic("That should not happen")
			}

			view := NewViewTask(todo, &m)

			return view.Update(nil)
		}

	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}
	return m.styles.BaseStyle.Render(m.table.View()) + "\n  " + m.help.View(m.keys) + "\n"
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

func (m Model) HandleDeleteTask() (table.Model, tea.Cmd) {

	cursor := m.table.Cursor()
	var NewCursor int
	if cursor == 0 {
		NewCursor = 0
	} else {
		NewCursor = cursor - 1
	}

	currentRow := m.table.SelectedRow()
	todoID, err := strconv.ParseInt(currentRow[0], 10, 64)
	if err != nil {
		panic("That should not happen")
	}
	todo, err := m.store.GetById(todoID)
	if err != nil {
		panic("That should not happen")
	}

	err = m.store.Delete(todo.ID)
	if nil != err {
		panic("That shouldn't happen")
	}
	todos, _ := m.store.GetAll()
	var rows []table.Row
	for _, todo := range todos {
		rows = append(rows, common.StringArray(todo))
	}
	m.table.SetRows(rows)
	m.table.SetCursor(NewCursor)
	m.table.UpdateViewport()
	return m.table.Update(nil)
}
