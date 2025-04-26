package cmd

import (
	"fmt"

	"github.com/cannable/sshcm/pkg/cdb"
	"github.com/spf13/cobra"
)

func listDefaults() error {
	fmt.Println("Program default settings:")

	for i := range cdb.ValidDefaults {
		def := cdb.ValidDefaults[i]

		val, err := db.GetDefault(def)

		if err != nil {
			return err
		}

		fmt.Printf("%-10s: %s\n", def, val)
	}

	return nil
}

// defaultsCmd represents the defaults command
var defaultsCmd = &cobra.Command{
	Use:   "defaults",
	Short: "List program defaults",
	Long: `
List program defaults.`,
	Example: `
sshcm defaults`,
	Run: func(cmd *cobra.Command, args []string) {
		db = openDb()

		err := listDefaults()

		if err != nil {
			panic(err)
		}

		db.Close()
	},
}

func init() {
	rootCmd.AddCommand(defaultsCmd)

	// Command flags
}
