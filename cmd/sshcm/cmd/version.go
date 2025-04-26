package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print program version",
	Long: `
Print program version and exit.`,
	Example: `
sshcm version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("SSH Connection Manager", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
