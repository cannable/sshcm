package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:     "search",
	Short:   "Search for connections",
	Long:    `Search for connections.`,
	Aliases: []string{"f"},
	Args:    cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("search called for", args[0])
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Command flags
}
