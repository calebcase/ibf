package cmd

import (
	"encoding/json"
	"os"

	"github.com/calebcase/ibf/lib"
	"github.com/spf13/cobra"
)

var unionCmd = &cobra.Command{
	Use:   "union IBF IBF [IBF]",
	Short: "Union the second IBF with the first. If third is provided, write the result there. Otherwise overwrite the first.",
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

		ibfs[0].Union(ibfs[1])

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
	RootCmd.AddCommand(unionCmd)
}
