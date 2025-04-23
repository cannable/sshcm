package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/cannable/ssh-cm-go/pkg/cdb"
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var (
	setNewNickname string

	setCmd = &cobra.Command{
		Use:     "set",
		Short:   "Alter an existing connection",
		Long:    `Alter an existing connection.`,
		Aliases: []string{"s"},
		Run: func(cmd *cobra.Command, args []string) {
			db = openDb()

			// Create a new connection struct and start populating it
			err := setConnection(cmd.Flags())

			if err != nil {
				if errors.Is(err, cdb.ErrConnectionNotFound) {
					fmt.Fprintln(os.Stderr, "Connection not updated because it could not be found.")
					os.Exit(1)
				} else if errors.Is(err, cdb.ErrNicknameLetter) {
					fmt.Fprintln(os.Stderr, "Nickname must begin with a letter.")
					os.Exit(1)
				} else {
					panic(err)
				}
			}

			if debugMode {
				fmt.Println("Updated connection.")
			}

			db.Close()
		},
	}
)

func init() {
	rootCmd.AddCommand(setCmd)

	// Command flags
	attachCommonCnFlags(setCmd, true)
	setCmd.PersistentFlags().StringVarP(&setNewNickname, "rename", "r", "", "Set new nickname for connection.")

	setCmd.MarkFlagsOneRequired("id", "nickname")
	setCmd.MarkFlagsMutuallyExclusive("id", "nickname")
}
