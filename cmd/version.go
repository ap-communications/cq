package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show cq version",
	Long:  "Show cq version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(VERSION) //print version information
	},
}

func init() {

	USAGE := `Usage:
  cq version
`
	RootCmd.AddCommand(versionCmd)
	versionCmd.SetUsageTemplate(USAGE) //override Usage words
}
