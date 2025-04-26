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

		cmd.Flags().Visit(accSetCnFlags)
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
	addCmd.PersistentFlags().StringVarP(&cmdCnNickname, "nickname", "n", "", "Nickname for connection")
	addCmd.PersistentFlags().StringVar(&cmdCnHost, "host", "", "Connection hostname (or IP address)")
	addCmd.PersistentFlags().StringVarP(&cmdCnUser, "user", "u", "", "User name for connection")
	addCmd.PersistentFlags().StringVarP(&cmdCnDescription, "description", "d", "", "Short description of the connection")
	addCmd.PersistentFlags().StringVarP(&cmdCnArgs, "args", "a", "", "Arguments to pass to SSH command")
	addCmd.PersistentFlags().StringVar(&cmdCnIdentity, "identity", "", "SSH identity to use for connection (a la '-i')")
	addCmd.PersistentFlags().StringVarP(&cmdCnCommand, "command", "c", "", "SSH command to run")

	addCmd.MarkPersistentFlagRequired("nickname")
	addCmd.MarkPersistentFlagRequired("host")
}
