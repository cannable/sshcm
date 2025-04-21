package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:     "set",
	Short:   "Alter an existing connection",
	Long:    `Alter an existing connection.`,
	Aliases: []string{"s"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("set called")
	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	// Command flags
	attachCommonCnFlags(setCmd, true)

	setCmd.MarkFlagsOneRequired("id", "nickname")
	setCmd.MarkFlagsMutuallyExclusive("id", "nickname")
}
