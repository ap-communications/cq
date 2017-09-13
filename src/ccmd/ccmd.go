package ccmd

import (
	"fmt"
	"os"

	"ccmd/commons"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var RootCmd = &cobra.Command{
	Use:   "cq",
	Short: "cq is a tool for simple and fast control cloud environment.",
	Long:  "cq is a tool for simple and fast control cloud environment.\nhttps://github.com/ap-communications/cq",
	Run: func(cmd *cobra.Command, args []string) {
		if commons.Flags.VersionFlag {
			printVersion()
		}
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.Flags().BoolVarP(&commons.Flags.VersionFlag, "version", "v", false, "Show cq version")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".cq")
	viper.AddConfigPath("$HOME")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
