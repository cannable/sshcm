package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
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
			fmt.Println("Connecting to ", c)
		}

		// Get effective SSH sshCmd (binary)
		sshCmd, err := db.GetEffectiveValue(c.Binary.Value, "binary")

		if err != nil {
			panic(err)
		}

		// If the program default is empty, use 'ssh'
		if strings.Compare(sshCmd, "") == 0 {
			sshCmd = "ssh"
		}

		// Make sure ssh binary resolves in PATH
		execBin, err := exec.LookPath("ssh")

		if err != nil {
			panic(err)
		}

		var execArgs = []string{execBin}

		// Append arguments
		sshArgs, err := db.GetEffectiveValue(c.Args.Value, "args")

		if err != nil {
			panic(err)
		}

		if strings.Compare(sshArgs, "") != 0 {
			// TODO: This is probably really mangled and won't work.
			// Figure out a way to reconstitute flat arguments from the DB.
			execArgs = append(execArgs, sshArgs)
		}

		// Append identity
		identity, err := db.GetEffectiveValue(c.Identity.Value, "identity")

		if err != nil {
			panic(err)
		}

		if strings.Compare(identity, "") != 0 {
			execArgs = append(execArgs, "-i", identity)
		}

		// Host & user
		host := c.Host.Value
		user, err := db.GetEffectiveValue(c.User.Value, "user")

		if err != nil {
			panic(err)
		}

		if strings.Compare(user, "") != 0 {
			execArgs = append(execArgs, user+"@"+host)
		} else {
			execArgs = append(execArgs, host)
		}

		if debugMode {
			fmt.Println("connection details:")
			fmt.Printf("binary:   '%s'\n", sshCmd)
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
