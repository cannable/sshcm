package cmd

import (
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

		if !isValidIdOrNickname(args[0]) {
			return ErrNoIdOrNickname
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		db = openDb()

		err := connect(args[0], cmd.Flags())

		if err != nil {
			panic(err)
		}

		db.Close()

	},
}

func init() {
	rootCmd.AddCommand(connectCmd)

	// Command flags
	attachCommonCnFlags(connectCmd, true)
}
