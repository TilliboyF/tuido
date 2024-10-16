package tui

import (
	"fmt"

	"github.com/TilliboyF/tuido/common"
	"github.com/TilliboyF/tuido/types"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type Form struct {
	/* Control Flow */
	isInit bool
	done   bool
	isNew  bool

	/* Form fields */
	todo      types.Todo
	lg        *lipgloss.Renderer
	styles    *FormStyles
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
	form := Form{
		todo:      todo,
		mainModel: mainModel,
		isInit:    true,
		isNew:     isNew,
	}
	form.lg = lipgloss.DefaultRenderer()
	form.styles = NewDefaultFormStyles(form.lg)

	var (
		nameInput = huh.NewInput().
				Key("name").
				Title("Name").
				Value(&todo.Name)

		descInput = huh.NewText().
				Key("description").
				Title("Description").
				Value(&todo.Description).Lines(3)

		confirm = huh.NewConfirm().
			Key("done").
			Title("All done?").
			Validate(func(v bool) error {
				if !v {
					return fmt.Errorf("Welp, finish up then")
				}
				return nil
			}).
			Affirmative("Yep").
			Negative("Wait, no")

		statusSelect = huh.NewSelect[string]().
				Key("status").
				Title("Status").
				Options(huh.NewOptions("todo", "in progress", "done")...).
				Value(common.PointerTo(todo.Status.String()))
	)

	if isNew {
		form.form = huh.NewForm(
			huh.NewGroup(
				nameInput,
				descInput,
				confirm,
			),
		)
	} else {
		form.form = huh.NewForm(
			huh.NewGroup(
				nameInput,
				descInput,
				statusSelect,
				confirm,
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
		return f, tea.Sequence(f.form.NextField(), f.form.PrevField()) // fixing a anoying bug
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
