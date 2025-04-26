package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// importCmd represents the import command
var (
	importFmt  string
	importPath string

	importCmd = &cobra.Command{
		Use:     "import",
		Short:   "Import connections",
		Long:    `Import connections.`,
		Aliases: []string{"i"},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("import called")
		},
	}
)

func init() {
	rootCmd.AddCommand(importCmd)

	// Command flags
	importCmd.PersistentFlags().StringVar(&importFmt, "format", "csv", "Import format. Valid formats: csv or json.")
	importCmd.PersistentFlags().StringVarP(&importPath, "path", "f", "", "Import path.")

}
