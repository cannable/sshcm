package cmd

import (
	"os"
	"text/template"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var (
	listAll bool

	listCmd = &cobra.Command{
		Use:     "list",
		Short:   "List all connections",
		Long:    `List all connections.`,
		Aliases: []string{"l"},
		Run: func(cmd *cobra.Command, args []string) {
			db = openDb()

			// Get all connections
			cns, err := db.GetAll()

			if err != nil {
				panic(err)
			}

			// Assemble output template
			t := `{{ printf "%-4s" "ID" }} `
			t = t + `{{ printf "%-15s" "Nickname" }} `
			t = t + `{{ printf "%-10s" "User" }} `
			t = t + `{{ printf "%-15s" "Host" }} `
			t = t + `{{ printf "%-20s" "Description" }} `

			if listAll {
				t = t + `{{ printf "%-10s" "Args" }} `
				t = t + `{{ printf "%-10s" "Identity" }} `
				t = t + `{{ printf "%-10s" "Command" }} `
				t = t + `{{ printf "%-10s"  "Binary" }} `
			}

			t = t + "\n{{ range . }}"
			t = t + `{{ .Id.StringTrimmed 4 }} `
			t = t + `{{ .Nickname.StringTrimmed 15 }} `
			t = t + `{{ .User.StringTrimmed 10 }} `
			t = t + `{{ .Host.StringTrimmed 15 }} `
			t = t + `{{ .Description.StringTrimmed 20 }} `

			if listAll {
				t = t + `{{ .Args.StringTrimmed 10 }} `
				t = t + `{{ .Identity.StringTrimmed 10 }} `
				t = t + `{{ .Command.StringTrimmed 10 }} `
				t = t + `{{ .Binary.StringTrimmed 10 }} `
			}

			t = t + "\n{{ end }}"

			tmpl, err := template.New("connection").Parse(t)

			if err != nil {
				panic(err)
			}

			// Run templates
			err = tmpl.Execute(os.Stdout, cns)

			if err != nil {
				panic(err)
			}

			db.Close()
		},
	}
)

func init() {
	rootCmd.AddCommand(listCmd)

	// Command flags
	listCmd.PersistentFlags().BoolVarP(&listAll, "all", "a", false, "List all connection details (wide output).")
}
