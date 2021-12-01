package cmd

import (
	"fmt"
	"os"

	ibf "github.com/calebcase/ibf/lib"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list IBF",
	Short: "List available keys from the set.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var path = args[0]

		set, err := open(path)
		if err != nil {
			return err
		}

		leftEmpty := true
		for val, err := set.Pop(); err == nil; val, err = set.Pop() {
			if !cfg.suppressLeft {
				fmt.Printf("%s\n", string(val))
			}
		}
		if !cfg.suppressLeft {
			leftEmpty = set.IsEmpty()
		}

		rightEmpty := true
		set.Invert()
		for val, err := set.Pop(); err == nil; val, err = set.Pop() {
			if !cfg.suppressRight {
				fmt.Printf("%s\n", string(val))
			}
		}
		if !cfg.suppressRight {
			rightEmpty = set.IsEmpty()
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

			return ibf.ErrNoPureCell
		}

		return nil
	},
}

func init() {
	listCmd.Flags().BoolVarP(&cfg.suppressLeft, "left", "1", false, "Suppress values unique to left-side (positive count).")
	listCmd.Flags().BoolVarP(&cfg.suppressRight, "right", "2", false, "Suppress values unique to right-side (negative count).")

	RootCmd.AddCommand(listCmd)
}
