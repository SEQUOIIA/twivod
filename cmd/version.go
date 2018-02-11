package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const VERSION string = "0.9.8"

var versionCmd = &cobra.Command{
	Use: "version",
	Short: "Print the version number of twiVod",
	Long: "Print the version number of twiVod",
	Run: func(cmd *cobra.Command, args[]string) {
		fmt.Printf("twiVod v%s\n", VERSION)
	},
}