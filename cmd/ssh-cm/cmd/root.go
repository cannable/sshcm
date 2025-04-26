package cmd

import (
	"fmt"
	"os"
	"slices"

	"github.com/cannable/ssh-cm-go/pkg/cdb"
	"github.com/spf13/cobra"
)

var (
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
		Use:   "ssh-cm",
		Short: "An SSH connection manager written in Go",
		Long:  `A simple SSH manager, written in Go, that uses a Sqlite DB.`,
	}
)

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
		cdb.ErrDbVersionNotRecognized,
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.PersistentFlags().StringVar(&connDbFilePath, "db", "", "Path to connection DB file (ssh-cm.connections).")
	rootCmd.PersistentFlags().BoolVarP(&debugMode, "verbose", "v", false, "Verbose output")
}
