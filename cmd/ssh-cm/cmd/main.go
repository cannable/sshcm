package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/cannable/ssh-cm-go/pkg/cdb"
	"github.com/spf13/pflag"
)

const Version = "v0.9"

var db cdb.ConnectionDB

func accSetCnFlags(f *pflag.Flag) {
	if cdb.IsValidProperty(f.Name) {
		cmdCnSetFlags = append(cmdCnSetFlags, f.Name)
	}
}

func getCnByIdOrNickname(arg string) (cdb.Connection, error) {
	var c cdb.Connection

	// Get connection by ID or nickname
	if err := cdb.ValidateId(arg); err == nil {
		// Got a valid ID
		id, err := strconv.Atoi(arg)

		if err != nil {
			return c, err
		}

		if debugMode {
			fmt.Println(id, "is an id.")
		}

		// Get connection by id
		c, err = db.Get(int64(id))

		if err != nil {
			return c, err
		}
	} else if err := cdb.ValidateNickname(arg); err == nil {
		// Got a valid nickname
		nickname := arg

		if debugMode {
			fmt.Println(nickname, "is a nickname.")
		}

		// Get connection by nickname
		c, err = db.GetByProperty("nickname", nickname)

		if err != nil {
			return c, err
		}
	}

	return c, nil
}

func getDbPath() string {
	const dbFileName = "ssh-cm.connections"

	/*
		Paths checked in this order:
			User-specified (ex. via argument)
			~/.config/dbFileName
			[current executable path]/dbFileName
	*/

	// Immediately return the path the user supplied, if they passed one
	if strings.Compare(connDbFilePath, "") != 0 {
		return connDbFilePath
	}

	// Assemble fallback path
	exe, err := os.Executable()

	if err != nil {
		// Panicking here is less than elegant, but a failure here is bad
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

func printConnection(c *cdb.Connection, printHeader bool) {
	t := `{{ printf "%-12s" "ID" }}: {{ .Id.Value }}
{{ printf "%-12s" "Nickname" }}: {{ .Nickname.Value }}
{{ printf "%-12s" "User" }}: {{ .User.Value }}
{{ printf "%-12s" "Host" }}: {{ .Host.Value }}
{{ printf "%-12s" "Description" }}: {{ .Description.Value }}
{{ printf "%-12s" "Args" }}: {{ .Args.Value }}
{{ printf "%-12s" "Identity" }}: {{ .Identity.Value }}
{{ printf "%-12s" "Command" }}: {{ .Command.Value }}
{{ printf "%-12s"  "Binary" }}: {{ .Binary.Value }}
`
	tmpl, err := template.New("connection_record").Parse(t)

	if err != nil {
		panic(err)
	}

	if printHeader {
		fmt.Println("******************************")
		fmt.Println(c.Nickname)
		fmt.Println("******************************")
	}

	// Run templates
	err = tmpl.Execute(os.Stdout, c)

	if err != nil {
		panic(err)
	}
}
