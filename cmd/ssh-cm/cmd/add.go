package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:     "add",
	Short:   "Add a connection",
	Long:    `Add a connection`,
	Aliases: []string{"a"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add called")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Command flags
	attachCommonCnFlags(addCmd, false)
	addCmd.MarkPersistentFlagRequired("nickname")
	addCmd.MarkPersistentFlagRequired("host")
}
