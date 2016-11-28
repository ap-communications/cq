package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var skelCmd = &cobra.Command{
	Use:   "skel",
	Short: "A short description",
	Long:  "A long description",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("skel called")
	},
}

func init() {
	RootCmd.AddCommand(skelCmd)
}
