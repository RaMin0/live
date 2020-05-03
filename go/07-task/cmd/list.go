package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all of your incomplete tasks",
	Run: func(cmd *cobra.Command, _ []string) {
		tasks, err := ListTasks(false)
		if err != nil {
			exitf("%v\n", err)
		}

		if len(tasks) == 0 {
			fmt.Println("You don't have any incomplete tasks.")
			os.Exit(0)
		}

		fmt.Println("You have the following tasks:")
		for _, t := range tasks {
			fmt.Printf("%d. %s\n", t.ID, t.Details)
		}
	},
}
