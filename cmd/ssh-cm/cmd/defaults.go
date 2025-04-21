package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// defaultsCmd represents the defaults command
var defaultsCmd = &cobra.Command{
	Use:   "defaults",
	Short: "List program defaults",
	Long:  `List program defaults.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("defaults called")
	},
}

func init() {
	rootCmd.AddCommand(defaultsCmd)

	// Command flags
}
