package cmd

import (
	"fmt"

	"github.com/cannable/sshcm/pkg/cdb"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get { id | nickname }",
	Short: "Print existing connection settings",
	Long: `
Print connection settings.

A valid connection ID or nickname must be specified.`,
	Example: `
sshcm get asdf
sshcm g 42
`,
	Aliases: []string{"g"},
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

		// Look up connection
		c, err := db.GetByIdOrNickname(args[0])

		if err != nil {
			bail(err)
		}

		// Show user the connection settings
		printConnection(&c, false)
		fmt.Println("")

		db.Close()
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Command flags
}
