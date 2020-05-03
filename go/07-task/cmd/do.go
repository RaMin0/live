package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(doCmd)
}

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Mark a task on your TODO list as complete",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			exitf("Missing task ID\n")
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			exitf("%v\n", err)
		}
		task := Task{ID: id}
		if err := MarkTaskAsCompleted(&task); err != nil {
			exitf("%v\n", err)
		}
		fmt.Printf("You have completed the %q task.\n", task.Details)
	},
}
