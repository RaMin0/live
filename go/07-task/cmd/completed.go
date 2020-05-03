package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(completedCmd)
}

var completedCmd = &cobra.Command{
	Use:   "completed",
	Short: "List all of your completed tasks",
	Run: func(cmd *cobra.Command, _ []string) {
		tasks, err := ListTasks(true)
		if err != nil {
			exitf("%v\n", err)
		}

		if len(tasks) == 0 {
			fmt.Println("You don't have any completed tasks.")
			os.Exit(0)
		}

		fmt.Println("You have finished the following tasks today:")
		for _, t := range tasks {
			fmt.Printf("%d. %s\n", t.ID, t.Details)
		}
	},
}
