package cmd

import (
	"encoding/csv"
	"os"

	"github.com/cannable/sshcm/pkg/cdb"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var (
	importFmt  string
	importPath string

	importCmd = &cobra.Command{
		Use:   "import",
		Short: "Import connections",
		Long: `
Import connections.

Connections will be imported appending to the existing list. Imported
connections are issued a new ID. Where a nickname already exists, the connection
will be updated to reflect what is in the imported data.`,
		Aliases: []string{"x"},
		Run: func(cmd *cobra.Command, args []string) {

			if len(importPath) > 0 {
				// Read from file

				if _, err := os.Stat(importPath); err != nil {
					bail(ErrImportFileNotFound)
				}

				// Open file
				f, err := os.Open(importPath)

				if err != nil {
					bail(err)
				}

				// Export
				err = importConnections(f)

				if err != nil {
					bail(err)
				}

				defer f.Close()
			} else {
				// Read from stdin
				err := importConnections(os.Stdin)

				if err != nil {
					bail(err)
				}
			}

		},
	}
)

func getCSVColumnMappings(row []string) (map[string]int, error) {
	cols := make(map[string]int)

	for id, col := range row {
		cols[col] = id
	}

	for col, _ := range cols {
		if !cdb.IsValidProperty(col) {
			return cols, ErrImportCSVInvalidColumn
		}
	}

	return cols, nil
}

func importConnections(f *os.File) error {
	db = openDb()

	cols := make(map[string]int)

	switch importFmt {

	case "csv":
		r := csv.NewReader(f)

		// Read all records from CSV file
		records, err := r.ReadAll()

		if err != nil {
			return err
		}

		// Loop through each record and import
		for i, row := range records {

			if i == 0 {
				// Header row - determine column order
				cols, err = getCSVColumnMappings(row)

			}

			c := cdb.NewConnection()

			c.Host = row[cols["host"]]

		}
	case "json":
		/*
			out, err := json.Marshal(cns)

			f.Write(out)

			if err != nil {
				return err
			}
		*/
	}
	db.Close()
	return nil
}

func init() {
	rootCmd.AddCommand(importCmd)

	// Command flags
	importCmd.PersistentFlags().StringVar(&importFmt, "format", "csv", "Export format. Valid formats: csv or json.")
	importCmd.PersistentFlags().StringVarP(&importPath, "path", "f", "", "Export destination path.")

}
