package cmd

import (
	"os"

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
