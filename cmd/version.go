package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version of karma
var Version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "shows karmahub version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("karmahub version", Version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
