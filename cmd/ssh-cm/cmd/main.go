package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/cannable/ssh-cm-go/pkg/cdb"
	"github.com/spf13/cobra"
)

var db cdb.ConnectionDB

// attachCommonCnFlags helper function that adds connection flags to the passed command.
func attachCommonCnFlags(cmd *cobra.Command, addId bool) {
	cmd.PersistentFlags().StringVarP(&cmdCnNickname, "nickname", "n", "", "Nickname for connection")
	cmd.PersistentFlags().StringVar(&cmdCnHost, "host", "", "Connection hostname (or IP address)")
	cmd.PersistentFlags().StringVarP(&cmdCnUser, "user", "u", "", "User name for connection")
	cmd.PersistentFlags().StringVarP(&cmdCnDescription, "description", "d", "", "Short description of the connection")
	cmd.PersistentFlags().StringVarP(&cmdCnArgs, "args", "a", "", "Arguments to pass to SSH command")
	cmd.PersistentFlags().StringVar(&cmdCnIdentity, "identity", "", "SSH identity to use for connection (a la '-i')")
	cmd.PersistentFlags().StringVarP(&cmdCnCommand, "command", "c", "", "SSH command to run")

	if addId {
		cmd.PersistentFlags().Int64VarP(&cmdCnId, "id", "i", 0, "ID of connection")
	}
}

func getDbPath() string {
	const dbFileName = "ssh-cm.connections"

	/*
		Paths to check, in this order:
			~/.config/dbFileName
			[current executable path]/dbFileName
	*/

	// Assemble fallback path
	exe, err := os.Executable()

	if err != nil {
		panic(err)
	}

	cwd := filepath.Dir(exe)
	fallback := filepath.Join(cwd, dbFileName)

	// Assemble preferred path
	homePath, err := os.UserHomeDir()

	if err != nil {
		return fallback
	}

	return filepath.Join(homePath, "/.config/"+dbFileName)
}

func listConnections(cns []*cdb.Connection, wide bool) {
	// Assemble output template
	t := `{{ printf "%-4s" "ID" }} | `
	t = t + `{{ printf "%-20s" "Nickname" }} | `
	t = t + `{{ printf "%-10s" "User" }} | `
	t = t + `{{ printf "%-20s" "Description" }} | `

	if wide {
		t = t + `{{ printf "%-10s" "Args" }} | `
		t = t + `{{ printf "%-10s" "Identity" }} | `
		t = t + `{{ printf "%-10s" "Command" }} | `
		t = t + `{{printf "%-10s"  "Binary" }} | `
	}

	t = t + "\n{{ range . }}"
	t = t + `{{ .Id | printf "%-4d" }} | `
	t = t + `{{ .Nickname | printf "%-20s" }} | `
	t = t + `{{ .User | printf "%-10s" }} | `
	t = t + `{{ .Description | printf "%-20s" }} | `

	if wide {
		t = t + `{{ .Args | printf "%-10s" }} | `
		t = t + `{{ .Identity | printf "%-10s" }} | `
		t = t + `{{ .Command | printf "%-10s" }} | `
		t = t + `{{ .Binary | printf "%-10s" }} | `
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
}

func openDb() cdb.ConnectionDB {
	path := getDbPath()

	if debugMode {
		fmt.Printf("Connection file path: '%s'\n", path)
	}

	// See if calling Open will create a new DB file
	if _, err := os.Stat(path); err != nil {
		fmt.Printf("Connection file '%s' does not exist and will be created.\n", path)
	}

	db, err := cdb.Open(path)

	if err != nil {
		panic(err)
	}

	return db
}
