package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/cannable/ssh-cm-go/pkg/cdb"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:     "remove",
	Short:   "Remove connection",
	Long:    `Remove connection.`,
	Aliases: []string{"rm", "delete", "del"},
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}

		if !cdb.IsValidIdOrNickname(args[0]) {
			return ErrNoIdOrNickname
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		db = openDb()

		c, err := getCnByIdOrNickname(args[0])

		if err != nil {
			if errors.Is(err, cdb.ErrConnNoId) {
				fmt.Fprintln(os.Stderr, "ID does not exist.")
				os.Exit(1)
			} else if errors.Is(err, cdb.ErrConnNoNickname) {
				fmt.Fprintln(os.Stderr, "Nickname does not exist.")
				os.Exit(1)
			} else if errors.Is(err, cdb.ErrConnectionNotFound) {
				fmt.Fprintln(os.Stderr, "Connection not found.")
				os.Exit(1)
			}
			panic(err)
		}

		if debugMode {
			fmt.Println("Deleting connection", c)
		}

		// Delete connection
		err = c.Delete()

		if err != nil {
			panic(err)
		}

		if err != nil {
			if errors.Is(err, cdb.ErrConnNoId) {
				fmt.Fprintln(os.Stderr, "ID does not exist.")
				os.Exit(1)
			} else if errors.Is(err, cdb.ErrConnNoNickname) {
				fmt.Fprintln(os.Stderr, "Nickname does not exist.")
				os.Exit(1)
			}
			panic(err)
		}

		db.Close()
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	// Command flags
}
