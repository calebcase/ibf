package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/calebcase/ibf/lib"
	"github.com/spf13/cobra"
)

var popCmd = &cobra.Command{
	Use:   "pop IBF",
	Short: "Remove and print the first available key from the set.",
	Run: func(cmd *cobra.Command, args []string) {
		var path = args[0]

		file, err := os.Open(path)
		cannot(err)

		decoder := json.NewDecoder(file)

		ibf := ibf.NewEmptyIBF()
		err = decoder.Decode(ibf)
		cannot(err)
		file.Close()

		val, err := ibf.Pop()
		cannot(err)

		fmt.Printf("%s\n", string(val.Bytes()))

		file, err = os.Create(path)
		cannot(err)

		encoder := json.NewEncoder(file)

		err = encoder.Encode(&ibf)
		cannot(err)
	},
}

func init() {
	RootCmd.AddCommand(popCmd)
}
