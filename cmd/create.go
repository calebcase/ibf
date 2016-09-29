package cmd

import (
	"encoding/json"
	"math/rand"
	"os"
	"strconv"

	"github.com/calebcase/ibf/lib"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create PATH SIZE [SEED]",
	Short: "Create a new set. Optionally specify a seed for the hash parameters.",
	Run: func(cmd *cobra.Command, args []string) {
		var path = args[0]
		var seed int64 = 0

		size, err := strconv.ParseUint(args[1], 10, 64)
		cannot(err)

		if len(args) > 2 {
			seed, err = strconv.ParseInt(args[2], 10, 64)
			cannot(err)
		}

		file, err := os.Create(path)
		cannot(err)

		r := rand.New(rand.NewSource(seed))

		positioners := []*ibf.Hasher{
			ibf.NewHasher(uint64(r.Int63()), uint64(r.Int63())),
			ibf.NewHasher(uint64(r.Int63()), uint64(r.Int63())),
			ibf.NewHasher(uint64(r.Int63()), uint64(r.Int63())),
		}
		hasher := ibf.NewHasher(uint64(r.Int63()), uint64(r.Int63()))

		set := ibf.NewIBF(size, positioners, hasher)

		enc := json.NewEncoder(file)

		err = enc.Encode(&set)
		cannot(err)
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
}
