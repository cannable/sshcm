package cmd

import (
	"github.com/spf13/cobra"
)

// defaultsCmd represents the defaults command
var defaultsCmd = &cobra.Command{
	Use:   "defaults",
	Short: "List program defaults",
	Long:  `List program defaults.`,
	Run: func(cmd *cobra.Command, args []string) {
		db = openDb()

		err := listDefaults()

		if err != nil {
			panic(err)
		}

		db.Close()
	},
}

func init() {
	rootCmd.AddCommand(defaultsCmd)

	// Command flags
}
