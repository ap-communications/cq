package ccmd

import (
	"fmt"

	"github.com/ap-communications/cq/src/ccmd/commons"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show cq version",
	Long:  "Show cq version",
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
	versionCmd.SetUsageTemplate("Usage: cq version")
}

func printVersion() {
	fmt.Println("Cloud Query (cq)   version " + commons.VERSION)
}
