/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/TilliboyF/tuido/common"
	"github.com/TilliboyF/tuido/store"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all todos",
	Long:  `List all unfinished todos, use flag -a to see all todos`,
	RunE:  list,
}

func list(cmd *cobra.Command, args []string) error {

	all, err := cmd.Flags().GetBool("all")
	if err != nil {
		return err
	}

	data, err := store.NewSqliteTodoStore(common.DB_LOCATION)
	if nil != err {
		return err
	}
	todos, err := data.GetAll()
	if err != nil {
		return err
	}

	if all {
		common.PrintTodos(todos)
	} else {
		filteredTodos := []store.Todo{}
		for _, todo := range todos {
			if !todo.Done {
				filteredTodos = append(filteredTodos, todo)
			}
		}
		common.PrintTodos(filteredTodos)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolP("all", "a", false, "list all todos")
}
