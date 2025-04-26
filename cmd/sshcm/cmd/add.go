package cmd

import (
	"fmt"
	"os"

	"github.com/cannable/sshcm/pkg/cdb"
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

		// Validate nickname follows the correct convention
		err := cdb.ValidateNickname(cmdCnNickname)

		if err != nil {
			bail(err)
		}

		// Nicknames must be unique. See if this one exists.
		exists := db.ExistsByProperty("nickname", cmdCnNickname)

		if exists {
			fmt.Fprintf(os.Stderr, "Can't add '%s'. Nickname already exists.\n", cmdCnNickname)
			os.Exit(1)
		}

		c := cdb.NewConnection()

		c.Nickname.Value = cmdCnNickname
		c.Host.Value = cmdCnHost
		c.User.Value = cmdCnUser
		c.Description.Value = cmdCnDescription
		c.Args.Value = cmdCnArgs
		c.Identity.Value = cmdCnIdentity
		c.Command.Value = cmdCnCommand

		if debugMode {
			fmt.Println("Adding connection:")
			printConnection(&c, false)
		}

		id, err := db.Add(&c)

		if err != nil {
			bail(err)
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
