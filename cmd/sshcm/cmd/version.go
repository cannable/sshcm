package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version = "v1.2.0-alpha.3"

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
