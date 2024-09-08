package cmd

import (
	"github.com/TilliboyF/tuido/common"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:                   "delete <id>",
	Short:                 "Delete a todo",
	Long:                  `Delete a todo`,
	DisableFlagsInUseLine: true,
	PreRunE:               common.ArgsCheckFunc(1),
	RunE:                  todoHandler.HandleDelete,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
