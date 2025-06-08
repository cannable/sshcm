package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
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
		Long: `Import connections from standard input (default) or a file.

The import process will update existing connections and append new ones.

The default format is CSV. To use json, pass --format json.`,
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

// getCSVColumnMappings returns a map of CSV column names vs. offset, determined
// from the passed row slice of strings. This allows importConnections() to
// re-assemble connection properties from CSV columns in any positional order,
// as long as they are valid properties.
//
// Returns map[string]int, where:
//
//		string == column heading (from passed row)
//	  int == positional index of column within row string slice/CSV file
//
// If a heading is encountered that is not valid, an error will be returned.
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

// importConnections imports connections from the passed Reader.
//
// nil will be returned if the entire import operation succeeds.
//
// This func supports csv and json format. The format used is determed by the
// global variable inputFmt. As this is managed by cobra/pflag, the value is
// not validated.
//
// Various errors may be returned at any point during the import and are
// likely caused by an I/O failure.
//
// An error during import will likely stop the process between connections
// (the previous one that succeeded and the current one that failed). As such,
// it's worth being aware that a partial import is a likely failure mode.
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

			importNickname := row[cols["nickname"]]

			// See if the nickname exists. If it does, we'll start with the existing
			// connection and update it
			update := false

			exists, err := db.ExistsByProperty("nickname", importNickname)

			if err != nil {
				return err
			}

			if exists {
				update = true
				c, err = db.GetByProperty("nickname", importNickname)

				// If we found the connection by nickname but couldn't actually retrieve
				// it, something is really wrong
				if err != nil {
					return err
				}

				fmt.Printf("Updating existing connection '%s' (%d)...\n", importNickname, c.Id)
			} else {
				fmt.Printf("Importing new connection '%s'...\n", importNickname)
			}

			// Populate or update properties in new Connection
			c.Nickname = row[cols["nickname"]]
			c.Host = row[cols["host"]]
			c.User = row[cols["user"]]
			c.Description = row[cols["description"]]
			c.Args = row[cols["args"]]
			c.Identity = row[cols["identity"]]
			c.Command = row[cols["command"]]

			if update {
				// Run smoke test on connection properties
				err = c.Validate()

				if err != nil {
					return err
				}

				// Update connection
				err = c.Update()

				if err != nil {
					return err
				}
			} else {
				// Run smoke test on connection properties
				err = c.Validate()

				if err != nil {
					return err
				}

				// Add connection
				id, err := db.Add(&c)

				if err != nil {
					return err
				}

				fmt.Printf("Added new connection '%s' (%d).\n", importNickname, id)
			}
		}
	case "json":
		d := json.NewDecoder(f)

		// Read the opening bracket (because we should have an array)
		_, err := d.Token()

		if err != nil {
			return err
		}

		// Read the stream of data, decoding Connection JSON payloads as we go
		for d.More() {
			c := cdb.NewConnection()

			err := d.Decode(&c)

			if err != nil {
				return err
			}

			// See if the nickname exists. If it does, we'll do an update.
			exists, err := db.ExistsByProperty("nickname", c.Nickname)

			if err != nil {
				return err
			}

			if exists {
				newCn := c

				c, err = db.GetByProperty("nickname", c.Nickname)

				// If we found the connection by nickname but couldn't actually retrieve
				// it, something is really wrong
				if err != nil {
					return err
				}

				// Update existing connection properties with those from the decoded
				// json object.
				c.Nickname = newCn.Nickname
				c.Host = newCn.Host
				c.User = newCn.User
				c.Description = newCn.Description
				c.Args = newCn.Args
				c.Identity = newCn.Identity
				c.Command = newCn.Command

				fmt.Printf("Updating existing connection '%s' (%d)...\n", c.Nickname, c.Id)

				// Run smoke test on connection properties
				err = c.Validate()

				if err != nil {
					return err
				}

				// Update connection

				err = c.Update()

				if err != nil {
					return err
				}
			} else {
				fmt.Printf("Importing new connection '%s'...\n", c.Nickname)

				// Run smoke test on connection properties
				err = c.Validate()

				if err != nil {
					return err
				}

				// Add connection
				id, err := db.Add(&c)

				if err != nil {
					return err
				}

				fmt.Printf("Added new connection '%s' (%d).\n", c.Nickname, id)
			}
		}
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
