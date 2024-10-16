package tui

import (
	"fmt"
	"strings"

	"github.com/TilliboyF/tuido/types"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ViewTask struct {
	todo      types.Todo
	mainModel *Model
	styles    *FormStyles
}

func NewViewTask(todo types.Todo, mainModel *Model) ViewTask {
	lg := lipgloss.DefaultRenderer()
	return ViewTask{
		todo:      todo,
		mainModel: mainModel,
		styles:    NewDefaultFormStyles(lg),
	}
}

func (v ViewTask) Init() tea.Cmd {
	return nil
}

func (v ViewTask) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return v.mainModel.Update(nil)
		}

	}
	return v, nil
}

func (v ViewTask) View() string {
	s := v.styles

	var b strings.Builder

	fmt.Fprintf(&b, "Name: %s \n", v.todo.Name)
	fmt.Fprintf(&b, "Description: \n")
	fmt.Fprintf(&b, "%s \n", v.todo.Description)
	fmt.Fprintf(&b, "Status: %s", v.todo.Status.String())

	return s.Status.Margin(0, 1).Padding(1, 2).Width(48).Render(b.String()) + "\n\n"
}
