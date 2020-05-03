package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rmCmd)
}

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Delete a task from your TODO list",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			exitf("Missing task ID\n")
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			exitf("%v\n", err)
		}
		task := Task{ID: id}
		if err := DeleteTask(&task); err != nil {
			exitf("%v\n", err)
		}
		fmt.Printf("You have deleted the %q task.\n", task.Details)
	},
}
