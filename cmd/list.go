/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/TilliboyF/tuido/common"
	"github.com/TilliboyF/tuido/store"
	"github.com/fatih/color"
	"github.com/mergestat/timediff"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: list,
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

	headerFmt := color.New(color.Underline).SprintfFunc()
	// columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("ID", "Name", "CreatedAt", "Done")
	tbl.WithHeaderFormatter(headerFmt) //.WithFirstColumnFormatter(columnFmt)

	for _, todo := range todos {
		if all {
			checkbox := common.Unchecked
			if todo.Done {
				checkbox = common.Checked
			}
			tbl.AddRow(todo.ID, todo.Name, timediff.TimeDiff(todo.CreatedAt), checkbox)
		} else {
			if !todo.Done {
				tbl.AddRow(todo.ID, todo.Name, timediff.TimeDiff(todo.CreatedAt), common.Unchecked)
			}
		}
	}

	tbl.Print()

	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolP("all", "a", false, "list all todos")
}
