package cmd

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"

	"github.com/calebcase/ibf/lib"
	"github.com/spf13/cobra"
)

func indexed(data []byte) (int64, []byte) {
	return int64(binary.LittleEndian.Uint64(data[len(data)-8:])), data[:len(data)-9]
}

var commCmd = &cobra.Command{
	Use:   "comm IBF1 IBF2",
	Short: "Compare IBF1 and IBF2.",
	Run: func(cmd *cobra.Command, args []string) {
		// Load IBF1 and IBF2.
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

		// Subtract IBF2 from IBF1.
		ibfs[0].Subtract(ibfs[1])
		ibf := ibfs[0]

		// Produce the two-column output.
		leftEmpty := true
		for val, err := ibf.Pop(); err == nil; val, err = ibf.Pop() {
			if !cfg.suppressLeft {
				if cfg.blockIndex >= 0 {
					idx, bytes := indexed(val.Bytes())
					fmt.Printf("%d:%s\n", idx, string(bytes))
				} else {
					fmt.Printf("%s\n", string(val.Bytes()))
				}
			}
		}
		if !cfg.suppressLeft {
			leftEmpty = ibf.IsEmpty()
		}

		rightEmpty := true
		ibf.Invert()
		for val, err := ibf.Pop(); err == nil; val, err = ibf.Pop() {
			if !cfg.suppressRight {
				if cfg.blockIndex >= 0 {
					idx, bytes := indexed(val.Bytes())
					fmt.Printf("%s%d:%s\n", cfg.columnDelimiter, idx, string(bytes))
				} else {
					fmt.Printf("%s%s\n", cfg.columnDelimiter, string(val.Bytes()))
				}
			}
		}
		if !cfg.suppressRight {
			rightEmpty = ibf.IsEmpty()
		}

		// Incomplete listing?
		if !leftEmpty || !rightEmpty {
			// Which side was empty?
			side := ""
			switch {
			case leftEmpty && rightEmpty:
				side = "left and right"
			case leftEmpty && !rightEmpty:
				side = "left"
			case !leftEmpty && rightEmpty:
				side = "right"
			}

			fmt.Fprintf(os.Stderr, "Unable to list all elements (%s).\n", side)
			os.Exit(1)
		}
	},
}

func init() {
	commCmd.Flags().StringVarP(&cfg.columnDelimiter, "output-delimiter", "d", "\t", "Separate columns with STR.")

	commCmd.Flags().BoolVarP(&cfg.suppressLeft, "left", "1", false, "Suppress values unique to left-side (IBF1).")
	commCmd.Flags().BoolVarP(&cfg.suppressRight, "right", "2", false, "Suppress values unique to right-side (IBF2).")

	commCmd.Flags().Int64VarP(&cfg.blockIndex, "block-index", "i", -1, "Values are assumed to be prefixed with an int64 index.")
	commCmd.Flags().Lookup("block-index").NoOptDefVal = "0"

	RootCmd.AddCommand(commCmd)
}
