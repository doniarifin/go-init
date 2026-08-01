package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version can be overridden at build time:
// go build -ldflags "-X github.com/doniarifin/go-init/cmd.Version=x.y.z"
var Version = "0.2.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print go-init version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("go-init version " + Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
