package cmd

import (
	"github.com/spf13/cobra"
)

var aclCmd = &cobra.Command{
	Use:   "acl",
	Short: "Control access control list",
	Long:  "Control access control list",
}

func init() {
	RootCmd.AddCommand(aclCmd)
}
