package cmd

import (
	"encoding/csv"
	"encoding/json"
	"os"

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
			db = openDb()

			// Get all connections
			cns, err := db.GetAll()

			if err != nil {
				bail(err)
			}

			cmd.Flags().Visit(accSetCnFlags)

			switch exportFmt {
			case "csv":
				w := csv.NewWriter(os.Stdout)

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
					bail(err)
				}

				// Write output
				for _, c := range cns {
					err = c.WriteCSV(w)

					if err != nil {
						bail(err)
					}
				}

				w.Flush()
			case "json":
				out, err := json.Marshal(cns)

				os.Stdout.Write(out)

				if err != nil {
					bail(err)
				}
			}

			db.Close()
		},
	}
)

func init() {
	rootCmd.AddCommand(exportCmd)

	// Command flags
	exportCmd.PersistentFlags().StringVar(&exportFmt, "format", "csv", "Export format. Valid formats: csv or json.")
	exportCmd.PersistentFlags().StringVarP(&exportPath, "path", "f", "", "Export path.")

}
