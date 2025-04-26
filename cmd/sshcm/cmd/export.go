package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// exportCmd represents the export command
var (
	exportFmt  string
	exportPath string

	exportCmd = &cobra.Command{
		Use:     "export",
		Short:   "Export all connections",
		Long:    `Export all connections.`,
		Aliases: []string{"x"},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("export called")
		},
	}
)

func init() {
	rootCmd.AddCommand(exportCmd)

	// Command flags
	exportCmd.PersistentFlags().StringVar(&exportFmt, "format", "csv", "Export format. Valid formats: csv or json.")
	exportCmd.PersistentFlags().StringVarP(&exportPath, "path", "f", "", "Export path.")

}
