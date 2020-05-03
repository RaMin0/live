package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "task",
	Short: "task is a CLI for managing your TODOs.",
}

// func init() {
// 	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
// }

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func exitf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	os.Exit(1)
}
