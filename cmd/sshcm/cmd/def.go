package cmd

import (
	"fmt"

	"github.com/cannable/sshcm/pkg/cdb"
	"github.com/spf13/cobra"
)

// defCmd represents the def command
var defCmd = &cobra.Command{
	Use:   "def setting value",
	Short: "Set program default settings",
	Long: `
Set program default settings.`,
	Example: `
sshcm def user asdf
`,
	Aliases: []string{},
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(2)(cmd, args); err != nil {
			return err
		}

		if !cdb.IsValidDefault(args[0]) {
			return ErrInvalidDefault
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		db = openDb()

		setting := args[0]
		value := args[1]

		err := db.SetDefault(setting, value)

		if err != nil {
			panic(err)
		}

		fmt.Printf("Updated '%s' default setting to '%s'.\n", setting, value)

		err = listDefaults()

		if err != nil {
			panic(err)
		}

		db.Close()
	},
}

func init() {
	rootCmd.AddCommand(defCmd)

	// Command flags
}
