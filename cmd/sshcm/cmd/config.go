package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// configCmd represents the config sub-command
var configCmd = &cobra.Command{
	Use:     "configuration",
	Aliases: []string{"config"},
	Short:   "Configuration commands",
	Long:    `Configuration commands.`,
}

// configShowCmd represents the configironment exists sub-sub-command
var configShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"print"},
	Short:   "Show configuration",
	Long:    `Print configuration to standard output.`,
	Run: func(cmd *cobra.Command, args []string) {
		path := getDbPath()

		fmt.Println("Connection DB: ", path)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
}
