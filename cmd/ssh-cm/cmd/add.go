package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/cannable/ssh-cm-go/pkg/cdb"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:     "add",
	Short:   "Add a connection",
	Long:    `Add a connection`,
	Aliases: []string{"a"},
	Run: func(cmd *cobra.Command, args []string) {
		db = openDb()

		// Create a new connection struct and start populating it
		id, err := addConnection()

		if err != nil {
			if errors.Is(err, cdb.ErrNicknameLetter) {
				fmt.Fprintln(os.Stderr, "Nickname must begin with a letter.")
				os.Exit(1)
			} else if errors.Is(err, ErrNicknameExists) {
				fmt.Fprintf(os.Stderr, "Can't add '%s'. Nickname already exists.\n", cmdCnNickname)
				os.Exit(1)
			} else {
				panic(err)
			}
		}

		fmt.Printf("Added new connection with id %d.\n", id)

		db.Close()
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Command flags
	attachCommonCnFlags(addCmd, false)
	addCmd.MarkPersistentFlagRequired("nickname")
	addCmd.MarkPersistentFlagRequired("host")
}
