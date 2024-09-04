/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/TilliboyF/tuido/common"
	"github.com/TilliboyF/tuido/store"
	"github.com/spf13/cobra"
)

// completeCmd represents the complete command
var completeCmd = &cobra.Command{
	Use:                   "complete <id>",
	Short:                 "Complete a task by the given id",
	Long:                  `Complete a task by the given id`,
	DisableFlagsInUseLine: true,
	PreRunE:               common.ArgsCheckFunc(1),
	RunE:                  complete,
}

func complete(cmd *cobra.Command, args []string) error {
	idString := args[0]
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		return fmt.Errorf("Given <id> = %s is not a integer", idString)
	}
	data, err := store.NewSqliteTodoStore(common.DB_LOCATION)
	if err != nil {
		return err
	}
	todo, err := data.GetById(id)
	if err != nil {
		return err
	}
	err = data.Complete(id)
	if err != nil {
		return err
	}
	todo.Done = true
	common.PrintTodo(todo)
	return nil
}

func init() {
	rootCmd.AddCommand(completeCmd)
}
