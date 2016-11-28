package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var VERSION string = "cq version 0.8.1" //version info
var cfgFile string

var RootCmd = &cobra.Command{
	Use:   "cq",
	Short: "cq is a tool for simple and fast control cloud environment.",
	Long: `cq is a tool for simple and fast control cloud environment.
https://github.com/ap-communications/cq`,
	Run: func(cmd *cobra.Command, args []string) {
		if listFlag.VersionFlag {
			fmt.Println(VERSION)
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
	RootCmd.Flags().BoolVarP(&listFlag.VersionFlag, "version", "v", false, "Show cq version") // --version -v

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	//	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cq.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".cq")   // name of config file (without extension)
	viper.AddConfigPath("$HOME") // adding home directory as first search path
	viper.AutomaticEnv()         // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
