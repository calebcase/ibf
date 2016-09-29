package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var echo string

var RootCmd = &cobra.Command{
	Use:   "set COMMAND [OPTIONS]",
	Short: "set CLI tool",
	Long:  "A CLI tool for managing set membership.",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	// Reads in config file and ENV variables.
	cobra.OnInitialize(initConfig)

	// Global configuration settings.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.set.yaml)")
}

func initConfig() {
	// Enable ability to specify config file via flag.
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	// Config file is named '.set.<ext>', loaded from the users home
	// directory, and overridden environment variables.
	viper.SetConfigName(".set")
	viper.AddConfigPath("$HOME")
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
