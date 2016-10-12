package cmd

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/calebcase/ibf/lib"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var insertCmd = &cobra.Command{
	Use:   "insert IBF [KEY]",
	Short: "Insert the key into the set. If key isn't provided, they will be read from stdin one per line.",
	Run: func(cmd *cobra.Command, args []string) {
		var path = args[0]

		// Should we echo our input?
		echoed := false
		if strings.Compare(cfg.echo, "true") == 0 {
			echoed = true
		} else if strings.Compare(cfg.echo, "false") == 0 {
			echoed = false
		} else if strings.Compare(cfg.echo, "auto") == 0 {
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
			ibf.Insert(val)
		} else {
			scanner := bufio.NewScanner(os.Stdin)

			if cfg.blockSize >= 0 {
				scanBlock := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
					if atEOF && len(data) == 0 {
						// At EOF and no more data to send.
						return 0, nil, nil
					}

					if len(data) == cfg.blockSize {
						// We have a complete block to send.
						return len(data), data, nil
					}

					if atEOF {
						// Send partial block.
						return len(data), data, nil
					}

					// Request more data.
					return 0, nil, nil
				}

				buf := make([]byte, cfg.blockSize)
				scanner.Buffer(buf, cfg.blockSize)
				scanner.Split(scanBlock)
			}

			count := -1

			for scanner.Scan() {
				bytes := scanner.Bytes()
				count += 1
				val := new(big.Int)

				if cfg.blockSize >= 0 && cfg.blockIndex >= 0 {
					idx := make([]byte, 8)
					binary.LittleEndian.PutUint64(idx, uint64(count))
					bytes = append(bytes, 1)
					bytes = append(bytes, idx...)
				}
				val.SetBytes(bytes)
				ibf.Insert(val)

				if echoed {
					fmt.Printf("%s\n", string(bytes))
				}
			}

			cannot(scanner.Err())
		}

		file, err = os.Create(path)
		cannot(err)

		encoder := json.NewEncoder(file)

		err = encoder.Encode(&ibf)
		cannot(err)
	},
}

func init() {
	insertCmd.Flags().StringVarP(&cfg.echo, "echo", "e", "auto", "Echo the values from stdin on stdout.")

	insertCmd.Flags().IntVarP(&cfg.blockSize, "block-size", "b", -1, "Set the block size for input parsing.")
	insertCmd.Flags().Lookup("block-size").NoOptDefVal = "4096"

	insertCmd.Flags().Int64VarP(&cfg.blockIndex, "block-index", "i", -1, "Suffix each block with an int64 index (starting at the provided value).")
	insertCmd.Flags().Lookup("block-index").NoOptDefVal = "0"

	RootCmd.AddCommand(insertCmd)
}
