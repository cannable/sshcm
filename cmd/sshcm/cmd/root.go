package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/cannable/sshcm/pkg/cdb"
	"github.com/cannable/sshcm/pkg/misc"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	db               cdb.ConnectionDB
	connDbFilePath   string
	debugMode        bool
	cmdCnNickname    string
	cmdCnHost        string
	cmdCnUser        string
	cmdCnDescription string
	cmdCnArgs        string
	cmdCnIdentity    string
	cmdCnCommand     string
	cmdCnSetFlags    []string

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "sshcm",
		Short: "An SSH connection manager written in Go",
		Long:  `A simple SSH manager, written in Go, that uses a Sqlite DB.`,
	}
)

// accSetCnFlags accumulates passed pflag names and stores them in the global
// slice cmdCnSetFlags.
func accSetCnFlags(f *pflag.Flag) {
	if cdb.IsValidProperty(f.Name) {
		cmdCnSetFlags = append(cmdCnSetFlags, f.Name)
	}
}

// bail reports somewhat-expected errors to the user in a "friendly" way.
// If the passed error is known and originates from the cdb module, this
// function will print the error to stderr and exit(1).
// If the error was not known, the program will panic.
func bail(err error) {
	minorErrors := []error{
		cdb.ErrConnNoDb,
		cdb.ErrConnNoId,
		cdb.ErrConnNoNickname,
		cdb.ErrConnectionNotFound,
		cdb.ErrDuplicateNickname,
		cdb.ErrIdNotExist,
		cdb.ErrInvalidConnectionProperty,
		cdb.ErrInvalidDefault,
		cdb.ErrInvalidId,
		cdb.ErrNicknameLetter,
		cdb.ErrPropertyInvalid,
		cdb.ErrSchemaVerInvalid,
	}

	if slices.Contains(minorErrors, err) && !debugMode {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	panic(err)
}

// getDbPath returns the path to the connection database.
//
// Paths checked in this order:
//
//	User-specified (ex. via argument)
//	~/.config/dbFileName
//	[current executable path]/dbFileName
func getDbPath() string {
	const dbFileName = "ssh-cm.connections"

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

// listConnections prints the passed connections in list format to stdout.
//
// wide controls whether all connection property columns are printed or a subset.
func listConnections(cns []*cdb.Connection, wide bool) {

	// Print header
	if wide {
		// Long header
		_, err := fmt.Fprintf(os.Stdout, "%s %s %s %s %s %s %s %s\n",
			misc.StringTrimmer("ID", cdb.ListViewColumnWidths["id"]),
			misc.StringTrimmer("Nickname", cdb.ListViewColumnWidths["nickname"]),
			misc.StringTrimmer("User", cdb.ListViewColumnWidths["user"]),
			misc.StringTrimmer("Host", cdb.ListViewColumnWidths["host"]),
			misc.StringTrimmer("Description", cdb.ListViewColumnWidths["description"]),
			misc.StringTrimmer("Args", cdb.ListViewColumnWidths["args"]),
			misc.StringTrimmer("Identity", cdb.ListViewColumnWidths["identity"]),
			misc.StringTrimmer("Command", cdb.ListViewColumnWidths["command"]),
		)

		if err != nil {
			bail(err)
		}
	} else {
		// Short header
		_, err := fmt.Fprintf(os.Stdout, "%s %s %s %s %s\n",
			misc.StringTrimmer("ID", cdb.ListViewColumnWidths["id"]),
			misc.StringTrimmer("Nickname", cdb.ListViewColumnWidths["nickname"]),
			misc.StringTrimmer("User", cdb.ListViewColumnWidths["user"]),
			misc.StringTrimmer("Host", cdb.ListViewColumnWidths["host"]),
			misc.StringTrimmer("Description", cdb.ListViewColumnWidths["description"]),
		)

		if err != nil {
			bail(err)
		}
	}

	for _, c := range cns {
		if wide {
			err := c.WriteLineLong(os.Stderr)

			if err != nil {
				bail(err)
			}
		} else {
			err := c.WriteLineShort(os.Stderr)

			if err != nil {
				bail(err)
			}

		}

	}

}

// openDb provides a simple wrapper around cdb.Open(). It calls getDbPath, then
// checks whether the path exists or not. If the connection DB file does not
// exist, it will print a message to stdout informing the user that one will
// be created. It then calls cdb.Open().
func openDb() cdb.ConnectionDB {
	path := getDbPath()

	if debugMode {
		fmt.Printf("Connection file path: '%s'\n", path)
	}

	// See if calling Open will create a new DB file
	create := false
	if _, err := os.Stat(path); err != nil {
		fmt.Printf("Connection file '%s' does not exist and will be created.\n", path)
		create = true
	}

	db, err := cdb.Connect("sqlite", path)

	if err != nil {
		panic(err)
	}

	// Create tables, if we need to. If not, see if upgrade is needed
	if create {
		err = db.InitializeDb(cdb.SchemaVersion)

		if err != nil {
			panic(err)
		}
	} else {
		// Can we use the DB?
		err = db.CheckDbHealth()

		if err != nil {
			switch err {
			case cdb.ErrSchemaUpgradeNeeded:
				// TODO: Do schema upgrade, if needed

			case cdb.ErrSchemaTooNew:
				// The schema version being too new for the tool is not a catastrophic
				// error.
				bail(err)
			default:
				panic(err)
			}
		}
	}

	return db
}

// printConnection writes connection properties to stdout in a multi-line
// record-format.
//
// printHeader controls whether a 'fancy' header string is printed.
func printConnection(c *cdb.Connection, printHeader bool) {
	if printHeader {
		fmt.Println("******************************")
		fmt.Println(c.Nickname)
		fmt.Println("******************************")
	}

	err := c.WriteRecordLong(os.Stdout)

	if err != nil {
		panic(err)
	}

}

// Execute adds all child commands to the root command and sets flags
// appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&connDbFilePath, "db", "", "Path to connection DB file (ssh-cm.connections).")
	rootCmd.PersistentFlags().BoolVarP(&debugMode, "verbose", "v", false, "Verbose output")
}
