package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfg struct {
	cfgFile         string
	echo            string
	suppressLeft    bool
	suppressRight   bool
	columnDelimiter string
	blockSize       int
	blockIndex      int64
}

var RootCmd = &cobra.Command{
	Use:   "ibf COMMAND [OPTIONS]",
	Short: "IBF CLI tool",
	Long:  "A CLI tool for managing Invertible Bloom Filters (IBF).",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Reads in config file and ENV variables.
	cobra.OnInitialize(initConfig)

	// Global configuration settings.
	RootCmd.PersistentFlags().StringVar(&cfg.cfgFile, "config", "", "config file (default is $HOME/.set.yaml)")
}

func initConfig() {
	// Enable ability to specify config file via flag.
	if cfg.cfgFile != "" {
		viper.SetConfigFile(cfg.cfgFile)
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
