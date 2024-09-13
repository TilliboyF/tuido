package handler

import (
	"fmt"
	"strconv"
	"time"

	"github.com/TilliboyF/tuido/common"
	"github.com/TilliboyF/tuido/db"
	"github.com/TilliboyF/tuido/tui"
	"github.com/TilliboyF/tuido/types"
	"github.com/spf13/cobra"

	tea "github.com/charmbracelet/bubbletea"
)

type TodoHandler struct {
	store *db.SqliteTodoStore
}

func NewTodoHandler(inMemoryDB bool) (*TodoHandler, error) {
	if store, _, err := db.NewSqliteTodoStore(false, inMemoryDB); err != nil {
		return nil, err
	} else {
		return &TodoHandler{
			store: store,
		}, nil
	}
}

func (h *TodoHandler) HandleAddTodo(cmd *cobra.Command, args []string) error {
	todo := types.Todo{
		Name: args[0],
	}
	err := h.store.Add(&todo)
	if err != nil {
		return err
	}
	todo.CreatedAt = time.Now()
	table := common.TableStringFromTodo(todo)
	cmd.Print(table)

	return nil
}

func (h *TodoHandler) HandleList(cmd *cobra.Command, args []string) error {

	all, err := cmd.Flags().GetBool("all")
	if err != nil {
		return err
	}

	todos, err := h.store.GetAll()
	if err != nil {
		return err
	}

	if all {
		table := common.TableStringFromTodos(todos)
		cmd.Print(table)
	} else {
		filteredTodos := []types.Todo{}
		for _, todo := range todos {
			if todo.Status == types.TODO {
				filteredTodos = append(filteredTodos, todo)
			}
		}
		table := common.TableStringFromTodos(filteredTodos)
		cmd.Print(table)
	}

	return nil
}

func (h *TodoHandler) HandleDelete(cmd *cobra.Command, args []string) error {
	idString := args[0]
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		return fmt.Errorf("Given <id> = %s is not a integer", idString)
	}
	_, err = h.store.GetById(id)
	if err != nil {
		return err
	}
	err = h.store.Delete(id)
	if err != nil {
		return err
	}
	cmd.Printf("Task id=%d deleted!\n", id)
	return nil
}

func (h *TodoHandler) HandleComplete(cmd *cobra.Command, args []string) error {
	idString := args[0]
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		return fmt.Errorf("Given <id> = %s is not a integer", idString)
	}
	todo, err := h.store.GetById(id)
	if err != nil {
		return err
	}
	err = h.store.Complete(id)
	if err != nil {
		return err
	}
	todo.Status = types.DONE
	table := common.TableStringFromTodo(todo)
	cmd.Print(table)
	return nil
}

func (h *TodoHandler) HandleMain(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		// TODO Open tui

		model := tui.NewModel(h.store)

		if _, err := tea.NewProgram(model).Run(); err != nil {
			return err
		}
	}
	return nil
}
