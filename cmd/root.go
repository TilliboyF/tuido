/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"

	"github.com/TilliboyF/tuido/handler"
	"github.com/spf13/cobra"
)

var TodoHandler *handler.TodoHandler

var rootCmd = &cobra.Command{
	Use:   "tuido",
	Short: "Todo list app",
	Long:  `Application to store and handle your todos in one place!`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	handler, err := handler.NewTodoHandler()
	if err != nil {
		log.Fatal("Error init: ", err)
	}
	TodoHandler = handler
}
