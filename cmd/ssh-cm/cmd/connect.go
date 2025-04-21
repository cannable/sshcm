package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:     "connect",
	Short:   "Start a connection",
	Long:    `Start a connection.`,
	Aliases: []string{"c"},
	Args:    cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("connect called for", args[0])
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)

	// Command flags
	attachCommonCnFlags(connectCmd, true)
}
