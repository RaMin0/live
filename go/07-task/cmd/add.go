package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task to your TODO list",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			exitf("Missing task details\n")
		}

		task := Task{Details: strings.Join(args, " ")}

		if err := CreateTask(&task); err != nil {
			exitf("%v\n", err)
		}
		fmt.Printf("Added %q to your task list.\n", task.Details)
	},
}
