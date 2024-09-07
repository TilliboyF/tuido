package cmd

import (
	"github.com/TilliboyF/tuido/common"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:                   "add <name>",
	Short:                 "adding a new todo",
	Long:                  `adding a new todo`,
	DisableFlagsInUseLine: true,
	PreRunE:               common.ArgsCheckFunc(1),
	RunE:                  TodoHandler.HandleAddTodo,
}

func init() {
	rootCmd.AddCommand(addCmd)
}
