package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// defaultsCmd represents the defaults command
var defaultsCmd = &cobra.Command{
	Use:   "defaults",
	Short: "List program defaults",
	Long:  `List program defaults.`,
	Run: func(cmd *cobra.Command, args []string) {
		db = openDb()

		fmt.Println("Program default settings:")

		defaults := [5]string{
			"binary",
			"user",
			"args",
			"identity",
			"command",
		}

		for i := range defaults {
			def := defaults[i]

			val, err := db.GetDefault(def)

			if err != nil {
				panic(err)
			}

			fmt.Printf("%-10s: %s\n", def, val)
		}

		db.Close()
	},
}

func init() {
	rootCmd.AddCommand(defaultsCmd)

	// Command flags
}
