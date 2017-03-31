package cmd

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/calebcase/ibf/lib"
	"github.com/spf13/cobra"
)

var mergeCmd = &cobra.Command{
	Use:   "merge IBF IBF [IBF]",
	Short: "Merge the second IBF into the first. If third is provided, write the result there. Otherwise overwrite the first. The difference between the first and second must be small enough to be completely listed otherwise merging is not possible.",
	Run: func(cmd *cobra.Command, args []string) {
		paths := args
		ibfs := [2]ibf.IBFer{}

		for i, path := range paths {
			if i > 1 {
				break
			}

			file, err := os.Open(path)
			cannot(err)

			decoder := json.NewDecoder(file)

			ibf := ibf.NewEmptyIBF()
			err = decoder.Decode(ibf)
			cannot(err)
			ibfs[i] = ibf

			file.Close()
		}

		// Attempt to remove the elements in first from second and then
		// list the remainder. For each of the remainder, insert into
		// first.
		ibfs[1].Subtract(ibfs[0])
		for val, err := ibfs[1].Pop(); err == nil; val, err = ibfs[1].Pop() {
			ibfs[0].Insert(val)
		}
		if !ibfs[1].IsEmpty() {
			cannot(errors.New("More elements in the set, but unable to retrieve."))
		}

		var output string

		if len(args) == 2 {
			output = paths[0]
		} else {
			output = paths[2]
		}

		file, err := os.Create(output)
		cannot(err)

		encoder := json.NewEncoder(file)

		err = encoder.Encode(ibfs[0])
		cannot(err)
	},
}

func init() {
	RootCmd.AddCommand(mergeCmd)
}
