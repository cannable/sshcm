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
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(1)(cmd, args); err != nil {
				return err
			}

			if !isValidIdOrNickname(args[0]) {
				return ErrNoIdOrNickname
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			db = openDb()

			cmd.Flags().Visit(accSetCnFlags)
			err := setConnection(args[0])

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
	setCmd.PersistentFlags().StringVarP(&cmdCnNickname, "nickname", "n", "", "Nickname for connection")
	setCmd.PersistentFlags().StringVar(&cmdCnHost, "host", "", "Connection hostname (or IP address)")
	setCmd.PersistentFlags().StringVarP(&cmdCnUser, "user", "u", "", "User name for connection")
	setCmd.PersistentFlags().StringVarP(&cmdCnDescription, "description", "d", "", "Short description of the connection")
	setCmd.PersistentFlags().StringVarP(&cmdCnArgs, "args", "a", "", "Arguments to pass to SSH command")
	setCmd.PersistentFlags().StringVar(&cmdCnIdentity, "identity", "", "SSH identity to use for connection (a la '-i')")
	setCmd.PersistentFlags().StringVarP(&cmdCnCommand, "command", "c", "", "SSH command to run")
}
