package cmd

import (
	"github.com/TilliboyF/tuido/common"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all todos",
	Long:    `List all unfinished todos, use flag -a to see all todos`,
	PreRunE: common.ArgsCheckFunc(0),
	RunE:    todoHandler.HandleList,
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolP("all", "a", false, "list all todos")
}
