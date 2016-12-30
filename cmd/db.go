package cmd

import (
	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Control manged DBs",
	Long:  "Control manged DBs",
}

func init() {
	RootCmd.AddCommand(dbCmd)
}
