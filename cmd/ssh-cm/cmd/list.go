package cmd

import (
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var (
	listAll bool

	listCmd = &cobra.Command{
		Use:     "list",
		Short:   "list all connections",
		Long:    `list all connections.`,
		Aliases: []string{"l"},
		Run: func(cmd *cobra.Command, args []string) {
			db = openDb()

			// Get all connections
			cns, err := db.GetAll()

			if err != nil {
				panic(err)
			}

			listConnections(cns, listAll)

			db.Close()
		},
	}
)

func init() {
	rootCmd.AddCommand(listCmd)

	// Command flags
	listCmd.PersistentFlags().BoolVarP(&listAll, "all", "a", false, "List all connection details (wide output).")
}
