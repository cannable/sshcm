package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"syscall"

	"github.com/cannable/sshcm/pkg/cdb"
	"github.com/spf13/cobra"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect { id | nickname }",
	Short: "Start a connection",
	Long: `
Start a connection.

A connection ID or nickname must be specified as the only positional argument.

Some connection settings (ex. command) can be overridden at runtime by passing flags.`,
	Example: `
sshcm connect something
sshcm c 22
sshcm c something --user=someone
`,
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
			c.User = cmdCnUser
		}

		if slices.Contains(cmdCnSetFlags, "args") {
			c.Args = cmdCnArgs
		}

		if slices.Contains(cmdCnSetFlags, "identity") {
			c.Identity = cmdCnIdentity
		}

		if slices.Contains(cmdCnSetFlags, "command") {
			c.Command = cmdCnCommand
		}

		if debugMode {
			fmt.Println("Connecting to ", c)
		}

		// Get effective SSH command
		sshCmd := c.Command
		if len(sshCmd) < 1 {
			sshCmd, err = db.GetDefault("command")

			if err != nil {
				bail(err)
			}
		}

		// If the program default is empty, use 'ssh'
		if len(sshCmd) < 1 {
			sshCmd = "ssh"
		}

		// Make sure ssh command resolves in PATH
		execBin, err := exec.LookPath(sshCmd)

		if err != nil {
			panic(err)
		}

		// Append arguments
		var execArgs = []string{execBin}

		sshArgs := c.Args
		if len(sshArgs) < 1 {
			sshArgs, err = db.GetDefault("args")

			if err != nil {
				panic(err)
			}
		}

		if len(sshArgs) > 0 {
			// TODO: This is probably really mangled and won't work.
			// Figure out a way to reconstitute flat arguments from the DB.
			execArgs = append(execArgs, sshArgs)
		}

		// Append identity
		identity := c.Identity

		if len(identity) < 1 {
			identity, err = db.GetDefault("identity")

			if err != nil {
				panic(err)
			}
		}

		if len(identity) > 0 {
			execArgs = append(execArgs, "-i", identity)
		}

		// Host
		host := c.Host

		// User
		user := c.User

		if len(user) < 1 {
			user, err = db.GetDefault("user")
			if err != nil {
				panic(err)
			}
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

		// We want to pass our environment to the new process
		execEnv := os.Environ()

		// Now's a good time to close the connection DB, since we're not going to
		// need it anymore
		db.Close()

		// Run the SSH command differently based on the OS on which we're running
		switch runtime.GOOS {
		case "windows":
			// On Windows, use os/exec to run the process
			exe := exec.Cmd{
				Path:   execBin,
				Args:   execArgs,
				Env:    execEnv,
				Stdin:  os.Stdin,
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			}

			err = exe.Run()

			if err != nil {
				fmt.Println("Error: ", err)
			}

		default:
			// On non-Windows systems, use syscall exec to replace the current process
			err = syscall.Exec(execBin, execArgs, execEnv)

			if err != nil {
				panic(err)
			}
		}
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
