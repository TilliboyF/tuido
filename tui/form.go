package tui

import (
	"fmt"

	"github.com/TilliboyF/tuido/types"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	red    = lipgloss.AdaptiveColor{Light: "#FE5F86", Dark: "#FE5F86"}
	indigo = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	green  = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
)

type Styles struct {
	Base,
	HeaderText,
	Status,
	StatusHeader,
	Highlight,
	ErrorHeaderText,
	Help lipgloss.Style
}

func NewStyles(lg *lipgloss.Renderer) *Styles {
	s := Styles{}
	s.Base = lg.NewStyle().
		Padding(1, 4, 0, 1)
	s.HeaderText = lg.NewStyle().
		Foreground(indigo).
		Bold(true).
		Padding(0, 1, 0, 2)
	s.Status = lg.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(indigo).
		PaddingLeft(1).
		MarginTop(1)
	s.StatusHeader = lg.NewStyle().
		Foreground(green).
		Bold(true)
	s.Highlight = lg.NewStyle().
		Foreground(lipgloss.Color("212"))
	s.ErrorHeaderText = s.HeaderText.
		Foreground(red)
	s.Help = lg.NewStyle().
		Foreground(lipgloss.Color("240"))
	return &s
}

type Form struct {
	/* Control Flow */
	isInit bool
	done   bool
	isNew  bool

	/* Form fields */
	todo      types.Todo
	lg        *lipgloss.Renderer
	styles    *Styles
	form      *huh.Form
	mainModel *Model
}

func NewForm(todo types.Todo, mainModel *Model) Form {
	return initForm(todo, false, mainModel)
}

func NewEmptyForm(mainModel *Model) Form {
	return initForm(types.Todo{}, true, mainModel)
}

func initForm(todo types.Todo, isNew bool, mainModel *Model) Form {
	form := Form{todo: todo, mainModel: mainModel, isInit: true, isNew: isNew}
	form.lg = lipgloss.DefaultRenderer()
	form.styles = NewStyles(form.lg)

	if isNew {
		form.form = huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("name").
					Title("Name").
					Value(&todo.Name),

				huh.NewText().
					Key("description").
					Title("Description").
					Value(&todo.Description),

				huh.NewConfirm().
					Key("done").
					Title("All done?").
					Validate(func(v bool) error {
						if !v {
							return fmt.Errorf("Welp, finish up then")
						}
						return nil
					}).
					Affirmative("Yep").
					Negative("Wait, no"),
			),
		)
	} else {
		form.form = huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("name").
					Title("Name").Value(&todo.Name),

				huh.NewText().
					Key("description").
					Title("Description").
					Value(&todo.Description),

				huh.NewSelect[string]().
					Key("status").
					Title("Status").
					Options(huh.NewOptions("todo", "in progress", "done")...),
				huh.NewConfirm().
					Key("done").
					Title("All done?").
					Validate(func(v bool) error {
						if !v {
							return fmt.Errorf("Welp, finish up then")
						}
						return nil
					}).
					Affirmative("Yep").
					Negative("Wait, no"),
			),
		)
	}

	return form
}

func (f Form) Init() tea.Cmd {
	return f.form.Init()
}

func (f Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			return f.mainModel.Update(nil)
		}
	}

	if f.isInit {

		f.isInit = false
		return f, tea.Sequence(f.form.NextField(), f.form.PrevField())
	}

	var cmds []tea.Cmd

	form, cmd := f.form.Update(msg)
	if newf, ok := form.(*huh.Form); ok {
		f.form = newf
		cmds = append(cmds, cmd)
	}
	if f.form.State == huh.StateCompleted {

		var status types.Status

		if !f.isNew {
			statusString := f.form.GetString("status")
			switch statusString {
			case "todo":
				status = types.TODO
			case "in progress":
				status = types.INPROGRESS
			case "done":
				status = types.DONE
			}
		} else {
			status = types.TODO
			f.todo.ID = -1
		}

		name := f.form.GetString("name")
		f.todo.Name = name
		desc := f.form.GetString("description")
		f.todo.Description = desc

		f.todo.Status = status

		return f.mainModel.Update(f.todo)
	}

	return f, tea.Sequence(cmds...)

}

func (f Form) View() string {
	return f.styles.Base.Render(f.form.View())
}
