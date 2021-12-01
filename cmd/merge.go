package cmd

import (
	"fmt"
	"os"

	ibf "github.com/calebcase/ibf/lib"
	"github.com/spf13/cobra"
)

var mergeCmd = &cobra.Command{
	Use:   "merge IBF IBF [IBF]",
	Short: "Merge the second IBF into the first. If third is provided, write the result there. Otherwise overwrite the first. The difference between the first and second must be small enough to be completely listed otherwise merging is not possible.",
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

		// Attempt to remove the elements in first from second and then
		// list the remainder. For each of the remainder, insert into
		// first.
		sets[1].Subtract(sets[0])
		for val, err := sets[1].Pop(); err == nil; val, err = sets[1].Pop() {
			sets[0].Insert(val)
		}
		if !sets[1].IsEmpty() {
			fmt.Fprintf(os.Stderr, "More elements in the set, but unable to retrieve.\n")

			return ibf.ErrNoPureCell
		}

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
	RootCmd.AddCommand(mergeCmd)
}
