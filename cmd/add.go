/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/TilliboyF/tuido/common"
	"github.com/TilliboyF/tuido/store"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:                   "add <name>",
	Short:                 "adding a new todo",
	Long:                  `adding a new todo`,
	DisableFlagsInUseLine: true,
	PreRunE:               common.ArgsCheckFunc(1),
	RunE:                  addTodo,
}

func addTodo(cmd *cobra.Command, args []string) error {
	data, err := store.NewSqliteTodoStore(common.DB_LOCATION)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		return fmt.Errorf("No argument for name given")
	}

	todo, err := data.Add(store.Todo{
		Name: args[0],
	})
	if err != nil {
		return err
	}
	todo.CreatedAt = time.Now()
	common.PrintTodo(todo)

	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)
}
