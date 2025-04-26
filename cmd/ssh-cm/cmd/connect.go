package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"syscall"

	"github.com/cannable/ssh-cm-go/pkg/cdb"
	"github.com/spf13/cobra"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:     "connect",
	Short:   "Start a connection",
	Long:    `Start a connection.`,
	Aliases: []string{"c"},
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

		cmd.Flags().Visit(accSetCnFlags)

		// Look up connection
		c, err := db.GetByIdOrNickname(args[0])

		if err != nil {
			bail(err)
		}

		// Override connection settings with user-supplied arguments
		if slices.Contains(cmdCnSetFlags, "user") {
			c.User.Value = cmdCnUser
		}

		if slices.Contains(cmdCnSetFlags, "args") {
			c.Args.Value = cmdCnArgs
		}

		if slices.Contains(cmdCnSetFlags, "identity") {
			c.Identity.Value = cmdCnIdentity
		}

		if slices.Contains(cmdCnSetFlags, "command") {
			c.Command.Value = cmdCnCommand
		}

		if debugMode {
			fmt.Println("Connecting to ", c)
		}

		// Get effective SSH command
		sshCmd, err := db.GetEffectiveValue(c.Command.Value, "command")

		if err != nil {
			bail(err)
		}

		// If the program default is empty, use 'ssh'
		if len(sshCmd) < 1 {
			sshCmd = "ssh"
		}

		// Make sure ssh binary resolves in PATH
		execBin, err := exec.LookPath(sshCmd)

		if err != nil {
			panic(err)
		}

		// Append arguments
		var execArgs = []string{execBin}
		sshArgs, err := db.GetEffectiveValue(c.Args.Value, "args")

		if err != nil {
			panic(err)
		}

		if len(sshArgs) > 0 {
			// TODO: This is probably really mangled and won't work.
			// Figure out a way to reconstitute flat arguments from the DB.
			execArgs = append(execArgs, sshArgs)
		}

		// Append identity
		identity, err := db.GetEffectiveValue(c.Identity.Value, "identity")

		if err != nil {
			panic(err)
		}

		if len(identity) > 0 {
			execArgs = append(execArgs, "-i", identity)
		}

		// Host & user
		host := c.Host.Value
		user, err := db.GetEffectiveValue(c.User.Value, "user")

		if err != nil {
			panic(err)
		}

		if len(user) > 0 {
			execArgs = append(execArgs, user+"@"+host)
		} else {
			execArgs = append(execArgs, host)
		}

		if debugMode {
			fmt.Println("connection details:")
			fmt.Printf("command:   '%s'\n", sshCmd)
			fmt.Printf("arguments:'%s'\n", execArgs)
		}

		// Connect
		execEnv := os.Environ()

		err = syscall.Exec(execBin, execArgs, execEnv)

		if err != nil {
			panic(err)
		}

		db.Close()
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)

	// Command flags
	connectCmd.PersistentFlags().StringVarP(&cmdCnUser, "user", "u", "", "User name for connection")
	connectCmd.PersistentFlags().StringVarP(&cmdCnArgs, "args", "a", "", "Arguments to pass to SSH command")
	connectCmd.PersistentFlags().StringVar(&cmdCnIdentity, "identity", "", "SSH identity to use for connection (a la '-i')")
	connectCmd.PersistentFlags().StringVarP(&cmdCnCommand, "command", "c", "", "SSH command to run")
}
