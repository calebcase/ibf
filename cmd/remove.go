package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/calebcase/ibf/lib"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var removeCmd = &cobra.Command{
	Use:   "remove IBF KEY",
	Short: "Remove the key into the set. If key isn't provided, they will be read from stdin one per line.",
	Run: func(cmd *cobra.Command, args []string) {
		var path = args[0]

		// Should we echo our input?
		echoed := false
		if strings.Compare(echo, "true") == 0 {
			echoed = true
		} else if strings.Compare(echo, "false") == 0 {
			echoed = false
		} else if strings.Compare(echo, "auto") == 0 {
			if !terminal.IsTerminal(int(os.Stdout.Fd())) {
				echoed = true
			}
		}

		file, err := os.Open(path)
		cannot(err)

		decoder := json.NewDecoder(file)

		ibf := ibf.NewEmptyIBF()
		err = decoder.Decode(ibf)
		cannot(err)
		file.Close()

		if len(args) == 2 {
			var key = args[1]

			val := new(big.Int)
			val.SetBytes([]byte(key))
			ibf.Remove(val)
		} else {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				bytes := scanner.Bytes()
				val := new(big.Int)
				val.SetBytes(bytes)
				ibf.Remove(val)

				if echoed {
					fmt.Printf("%s\n", string(bytes))
				}
			}
		}

		file, err = os.Create(path)
		cannot(err)

		encoder := json.NewEncoder(file)

		err = encoder.Encode(&ibf)
		cannot(err)
	},
}

func init() {
	removeCmd.Flags().StringVarP(&echo, "echo", "e", "auto", "Echo the values from stdin on stdout.")

	RootCmd.AddCommand(removeCmd)
}
