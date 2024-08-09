/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

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
	// Args: func(cmd *cobra.Command, args []string) error {
	// 	if len(args) < 1 {
	// 		return fmt.Errorf("requires a name argument")
	// 	}
	// 	return nil
	// },
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
		return nil
	},
	RunE: addTodo,
}

func addTodo(cmd *cobra.Command, args []string) error {
	// fmt.Println("add called")
	// test, _ := cmd.Flags().GetString("todo")
	// fmt.Println("flag: ", test)
	// fmt.Println("args: ", args)

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
	fmt.Println("New todo: ", todo)

	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)
	// addCmd.Flags().String("todo", "", "wvwvwrv")
}
