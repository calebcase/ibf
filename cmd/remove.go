package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var removeCmd = &cobra.Command{
	Use:   "remove IBF KEY",
	Short: "Remove the key into the set. If key isn't provided, they will be read from stdin one per line.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
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

		set, err := open(path)
		if err != nil {
			return err
		}

		if len(args) == 2 {
			set.Remove([]byte(args[1]))
		} else {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				bytes := scanner.Bytes()

				set.Remove(bytes)

				if echoed {
					fmt.Printf("%s\n", string(bytes))
				}
			}
			err = scanner.Err()
			if err != nil {
				return err
			}
		}

		return create(path, set)
	},
}

func init() {
	removeCmd.Flags().StringVarP(&cfg.echo, "echo", "e", "auto", "Echo the values from stdin on stdout.")

	RootCmd.AddCommand(removeCmd)
}
