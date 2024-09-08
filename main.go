package main

import (
	"fmt"
	"log"

	"github.com/TilliboyF/tuido/common"
	"github.com/TilliboyF/tuido/handler"
	"github.com/spf13/cobra"
)

var (
	version = "1.0.2"
)

func NewRootCmd(handler *handler.TodoHandler) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "tuido",
		Short: "Todo list app",
		Long:  `Application to store and handle your todos in one place!`,
	}
	rootCmd.AddCommand(NewAddCmd(handler))
	rootCmd.AddCommand(NewCompleteCmd(handler))
	rootCmd.AddCommand(NewListCmd(handler))
	rootCmd.AddCommand(NewDeleteCmd(handler))
	rootCmd.AddCommand(NewVersionCmd())

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	return rootCmd
}

func NewAddCmd(handler *handler.TodoHandler) *cobra.Command {
	var addCmd = &cobra.Command{
		Use:                   "add <name>",
		Short:                 "adding a new todo",
		Long:                  `adding a new todo`,
		DisableFlagsInUseLine: true,
		PreRunE:               common.ArgsCheckFunc(1),
		RunE:                  handler.HandleAddTodo,
	}
	return addCmd
}

func NewCompleteCmd(handler *handler.TodoHandler) *cobra.Command {
	completeCmd := &cobra.Command{
		Use:                   "complete <id>",
		Short:                 "Complete a task by the given id",
		Long:                  `Complete a task by the given id`,
		DisableFlagsInUseLine: true,
		PreRunE:               common.ArgsCheckFunc(1),
		RunE:                  handler.HandleComplete,
	}
	return completeCmd
}

func NewListCmd(handler *handler.TodoHandler) *cobra.Command {
	listCmd := &cobra.Command{
		Use:     "list",
		Short:   "List all todos",
		Long:    `List all unfinished todos, use flag -a to see all todos`,
		PreRunE: common.ArgsCheckFunc(0),
		RunE:    handler.HandleList,
	}
	listCmd.Flags().BoolP("all", "a", false, "list all todos")
	return listCmd
}

func NewDeleteCmd(handler *handler.TodoHandler) *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:                   "delete <id>",
		Short:                 "Delete a todo",
		Long:                  `Delete a todo`,
		DisableFlagsInUseLine: true,
		PreRunE:               common.ArgsCheckFunc(1),
		RunE:                  handler.HandleDelete,
	}
	return deleteCmd
}

func NewVersionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:                   "version",
		Short:                 "Get version of app",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("tuido app version %s\n", version)
		},
	}
	return versionCmd
}

func main() {
	handler, err := handler.NewTodoHandler("")
	if err != nil {
		log.Fatal(err)
	}
	rootCmd := NewRootCmd(handler)
	rootCmd.Execute()
}
