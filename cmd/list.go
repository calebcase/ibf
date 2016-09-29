package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/calebcase/ibf/lib"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list IBF",
	Short: "List available keys from the set.",
	Run: func(cmd *cobra.Command, args []string) {
		var path = args[0]

		file, err := os.Open(path)
		cannot(err)

		decoder := json.NewDecoder(file)

		ibf := &ibf.IBF{}
		err = decoder.Decode(ibf)
		cannot(err)
		file.Close()

		for val, err := ibf.Pop(); err == nil; val, err = ibf.Pop() {
			fmt.Printf("%s\n", string(val.Bytes()))
		}

		// Incomplete listing.
		if !ibf.IsEmpty() {
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
