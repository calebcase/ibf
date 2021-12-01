package cmd

import (
	"github.com/spf13/cobra"
)

var invertCmd = &cobra.Command{
	Use:   "invert IBF",
	Short: "Invert the counts of the IBF (multiply by -1).",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var path = args[0]

		set, err := open(path)
		if err != nil {
			return err
		}

		set.Invert()

		return create(path, set)
	},
}

func init() {
	RootCmd.AddCommand(invertCmd)
}
