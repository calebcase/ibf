package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var popCmd = &cobra.Command{
	Use:   "pop IBF",
	Short: "Remove and print the first available key from the set.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var path = args[0]

		set, err := open(path)
		if err != nil {
			return err
		}

		val, err := set.Pop()
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", string(val))

		return create(path, set)
	},
}

func init() {
	RootCmd.AddCommand(popCmd)
}
