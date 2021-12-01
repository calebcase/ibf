package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion SHELL",
	Short: "Output shell completion code for the given shell.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		shell := args[0]

		switch shell {
		case "bash":
			return RootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return RootCmd.GenZshCompletion(os.Stdout)
		case "powershell":
			return RootCmd.GenPowerShellCompletion(os.Stdout)
		}

		return errors.New("unknown shell")
	},
}

func init() {
	RootCmd.AddCommand(completionCmd)
}
