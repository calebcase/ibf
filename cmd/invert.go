package cmd

import (
	"encoding/json"
	"os"

	"github.com/calebcase/ibf/lib"
	"github.com/spf13/cobra"
)

var invertCmd = &cobra.Command{
	Use:   "invert IBF",
	Short: "Invert the counts of the IBF (multiply by -1).",
	Run: func(cmd *cobra.Command, args []string) {
		var path = args[0]

		file, err := os.Open(path)
		cannot(err)

		decoder := json.NewDecoder(file)

		ibf := ibf.NewEmptyIBF()
		err = decoder.Decode(ibf)
		cannot(err)
		file.Close()

		ibf.Invert()

		file, err = os.Create(path)
		cannot(err)

		encoder := json.NewEncoder(file)

		err = encoder.Encode(&ibf)
		cannot(err)
	},
}

func init() {
	RootCmd.AddCommand(invertCmd)
}
