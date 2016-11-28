package cmd

import (
	"github.com/spf13/cobra"
)

var vmCmd = &cobra.Command{
	Use:   "vm",
	Short: "Control virtual machines",
	Long:  "Control virtual machines",
}

func init() {
	RootCmd.AddCommand(vmCmd)
}
