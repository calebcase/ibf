package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion SHELL",
	Short: "Output shell completion code for the given shell.",
	Run: func(cmd *cobra.Command, args []string) {
		shell := args[0]

		switch shell {
		case "bash":
			err := RootCmd.GenBashCompletion(os.Stdout)
			cannot(err)
		default:
			fmt.Fprintf(os.Stderr, "Unknown shell.\n")
		}
	},
}

func init() {
	RootCmd.AddCommand(completionCmd)
}
