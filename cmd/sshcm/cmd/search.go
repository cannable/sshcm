package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for connections",
	Long: `Search for connections. Search is case-insensitive.

Search references the following properties:

- nickname
- host
- user
- description`,
	Aliases: []string{"f"},
	Args:    cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		if debugMode {
			fmt.Println("Searching for '", args[0]+"'")

		}

		db = openDb()

		// Get all connections
		cns, err := db.Search(args[0])

		if err != nil {
			panic(err)
		}

		listConnections(cns, listAll)

		db.Close()
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Command flags
	searchCmd.PersistentFlags().BoolVarP(&listAll, "all", "a", false, "List all connection details (wide output).")
}
