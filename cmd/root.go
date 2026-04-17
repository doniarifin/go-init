package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-init",
	Short: "Go project boilerplate generator",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to go-init :)")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
