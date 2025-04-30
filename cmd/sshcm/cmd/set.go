package cmd

import (
	"fmt"
	"slices"
	"strings"

	"github.com/cannable/sshcm/pkg/cdb"
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set { id | nickname }",
	Short: "Change connection settings",
	Long: `Change connection settings.
A valid ID or nickname must be specified.

A connection can be renamed by passing --nickname="new_nickname".
`,
	Example: `
sshcm set 42 --user="blarg"
sshcm s asdf --nickname fdsa`,
	Aliases: []string{"s"},
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

		oldNickname := args[0]

		cmd.Flags().Visit(accSetCnFlags)

		// Look up connection
		c, err := db.GetByIdOrNickname(oldNickname)

		if err != nil {
			bail(err)
		}

		// Show original values if in debug mode
		if debugMode {
			fmt.Println("Current connection settings:")
			printConnection(&c, false)
			fmt.Println("")
		}

		// Determine if we're renaming
		if slices.Contains(cmdCnSetFlags, "nickname") {
			// Validate nickname follows the correct convention
			err = cdb.ValidateNickname(cmdCnNickname)

			if err != nil {
				bail(err)
			}

			// See if the new nickname exists already.
			exists := db.ExistsByProperty("nickname", cmdCnNickname)

			if exists {
				// Bail if the nickname already exists and isn't the current connection
				if strings.Compare(cmdCnNickname, oldNickname) != 0 {
					bail(cdb.ErrDuplicateNickname)
				}
			}

			c.Nickname = cmdCnNickname
		}

		// Update hostname, if it was passed
		if slices.Contains(cmdCnSetFlags, "host") {
			c.Host = cmdCnHost
		}

		// Update user, if it was passed
		if slices.Contains(cmdCnSetFlags, "user") {
			c.User = cmdCnUser
		}

		// Update description, if it was passed
		if slices.Contains(cmdCnSetFlags, "description") {
			c.Description = cmdCnDescription
		}

		// Update args, if it was passed
		if slices.Contains(cmdCnSetFlags, "args") {
			c.Args = cmdCnArgs
		}

		// Update identity, if it was passed
		if slices.Contains(cmdCnSetFlags, "identity") {
			c.Identity = cmdCnIdentity
		}

		// Update command, if it was passed
		if slices.Contains(cmdCnSetFlags, "command") {
			c.Command = cmdCnCommand
		}

		// Run smoke test on connection properties
		err = c.Validate()

		if err != nil {
			bail(err)
		}

		// Update the connection
		err = c.Update()

		if err != nil {
			panic(err)
		}

		// Show user the updated connection settings
		fmt.Println("New connection settings:")
		printConnection(&c, false)
		fmt.Println("")

		db.Close()
	},
}

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
