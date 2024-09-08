package common

import (
	"os"

	"github.com/fatih/color"
	"github.com/mergestat/timediff"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"

	"github.com/TilliboyF/tuido/types"
)

const (
	DB_NAME   = "data.db"
	Unchecked = "☐"
	Checked   = "☑"
)

func PrintTodos(todos []types.Todo) {
	headerFmt := color.New(color.FgCyan, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("ID", "Name", "CreatedAt", "Done")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, todo := range todos {
		tbl.AddRow(todo.ID, todo.Name, timediff.TimeDiff(todo.CreatedAt), getCheckBox(todo.Done))
	}

	tbl.Print()

}

func PrintTodo(todo types.Todo) {
	headerFmt := color.New(color.FgCyan, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("ID", "Name", "CreatedAt", "Done")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	tbl.AddRow(todo.ID, todo.Name, timediff.TimeDiff(todo.CreatedAt), getCheckBox(todo.Done))

	tbl.Print()
}

func getCheckBox(done bool) string {
	checkbox := Unchecked
	c := color.FgRed
	if done {
		checkbox = Checked
		c = color.FgGreen
	}
	return color.New(c).Sprint(checkbox)
}

func ArgsCheckFunc(amount int) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != amount {
			cmd.Help()
			os.Exit(0)
		}
		return nil
	}
}
