package cmd

import (
	ibf "github.com/calebcase/ibf/lib"
	"github.com/spf13/cobra"
)

var unionCmd = &cobra.Command{
	Use:   "union IBF IBF [IBF]",
	Short: "Union the second IBF with the first. If third is provided, write the result there. Otherwise overwrite the first.",
	Args:  cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		paths := args
		sets := [2]*ibf.IBF{}

		for i, path := range paths {
			if i > 1 {
				break
			}

			sets[i], err = open(path)
			if err != nil {
				return err
			}
		}

		sets[0].Union(sets[1])

		var output string

		if len(args) == 2 {
			output = paths[0]
		} else {
			output = paths[2]
		}

		return create(output, sets[0])
	},
}

func init() {
	RootCmd.AddCommand(unionCmd)
}
