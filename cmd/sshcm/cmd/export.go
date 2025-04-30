package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// exportCmd represents the export command
var (
	exportFmt  string
	exportPath string

	exportCmd = &cobra.Command{
		Use:   "export",
		Short: "Export all connections",
		Long:  `Export all connections.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(exportPath) > 0 {
				// Write to file

				if _, err := os.Stat(exportPath); err == nil {
					// Print warning because the output file exists
					fmt.Fprintln(os.Stderr,
						"warning: export file exists and will be overwritten")
				}

				// Open file
				f, err := os.OpenFile(exportPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0640)

				if err != nil {
					bail(err)
				}

				// Export
				err = exportConnections(f)

				if err != nil {
					bail(err)
				}

				defer f.Close()
			} else {
				// Write to stdout
				err := exportConnections(os.Stdout)

				if err != nil {
					bail(err)
				}
			}

		},
	}
)

func exportConnections(f *os.File) error {
	db = openDb()

	// Get all connections
	cns, err := db.GetAll()

	if err != nil {
		bail(err)
	}

	switch exportFmt {

	case "csv":
		w := csv.NewWriter(f)

		// Write CSV file header
		header := []string{
			"id",
			"nickname",
			"user",
			"host",
			"description",
			"args",
			"identity",
			"command",
		}

		err := w.Write(header)

		if err != nil {
			return err
		}

		// Write output
		for _, c := range cns {
			err = c.WriteCSV(w)

			if err != nil {
				return err
			}
		}

		w.Flush()
	case "json":
		out, err := json.Marshal(cns)

		f.Write(out)

		if err != nil {
			return err
		}
	}

	db.Close()

	return nil
}

func init() {
	rootCmd.AddCommand(exportCmd)

	// Command flags
	exportCmd.PersistentFlags().StringVar(&exportFmt, "format", "csv", "Export format. Valid formats: csv or json.")
	exportCmd.PersistentFlags().StringVarP(&exportPath, "path", "f", "", "Export destination path.")

}
