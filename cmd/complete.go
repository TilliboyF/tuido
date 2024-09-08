package cmd

import (
	"github.com/TilliboyF/tuido/common"
	"github.com/spf13/cobra"
)

var completeCmd = &cobra.Command{
	Use:                   "complete <id>",
	Short:                 "Complete a task by the given id",
	Long:                  `Complete a task by the given id`,
	DisableFlagsInUseLine: true,
	PreRunE:               common.ArgsCheckFunc(1),
	RunE:                  todoHandler.HandleComplete,
}

func init() {
	rootCmd.AddCommand(completeCmd)
}
