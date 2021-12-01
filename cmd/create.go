package cmd

import (
	"strconv"

	ibf "github.com/calebcase/ibf/lib"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create PATH SIZE [SEED]",
	Short: "Create a new set. Optionally specify a seed for the hash parameters.",
	Args:  cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var path = args[0]
		var seed int64 = 0

		size, err := strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			return err
		}

		if len(args) > 2 {
			seed, err = strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return err
			}
		}

		set := ibf.NewIBF(size, seed)

		return create(path, set)
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
}
