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

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:                   "delete <id>",
	Short:                 "Delete a todo",
	Long:                  `Delete a todo`,
	DisableFlagsInUseLine: true,
	PreRunE:               common.ArgsCheckFunc(1),
	RunE:                  delete,
}

func delete(cmd *cobra.Command, args []string) error {
	idString := args[0]
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		return fmt.Errorf("Given <id> = %s is not a integer", idString)
	}
	data, err := store.NewSqliteTodoStore(common.DB_LOCATION)
	if err != nil {
		return err
	}
	_, err = data.GetById(id)
	if err != nil {
		return err
	}
	err = data.Delete(id)
	if err != nil {
		return err
	}
	fmt.Printf("Task id=%d deleted!\n", id)
	return nil
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
