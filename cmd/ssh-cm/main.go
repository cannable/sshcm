package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func helpBlurb() string {
	var b strings.Builder

	// Buiild file contents
	fmt.Fprintf(&b, `SSH Connection Manager
This is a simple SSH wrapper tool rewritten in Go by C.Annable.

`)

	return b.String()
}

func main() {
	homePath, err := os.UserHomeDir()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Home dir:", homePath)

	// Semi-common connection-related flags
	commonCnFlags := map[string]string{
		"args":        "Arguments to pass to SSH command",
		"command":     "SSH command to run",
		"description": "Short description of the connection",
		"host":        "Connection hostname (or IP address)",
		"identity":    "SSH identity to use for connection (a la '-i')",
		"nickname":    "Nickname for connection",
		"user":        "User name for connection",
	}

	// Flag definitions - Add connection
	addFlagSet := flag.NewFlagSet("add", flag.ExitOnError)
	addCmdCnSettings := make(map[string]*string)

	for f, s := range commonCnFlags {
		addCmdCnSettings[f] = addFlagSet.String(f, "", s)
	}

	// Flag definitions - Set connection property
	setFlagSet := flag.NewFlagSet("set", flag.ExitOnError)
	setCmdCnSettings := make(map[string]*string)
	setID := setFlagSet.Int("id", 0, "Connection id")

	for f, s := range commonCnFlags {
		setCmdCnSettings[f] = setFlagSet.String(f, "", s)
	}

	// Flag definitions - Export connections
	exportFlagSet := flag.NewFlagSet("export", flag.ExitOnError)
	exportFormat := exportFlagSet.String("format", "csv", "export content format")

	// Flag definitions - Import connections
	importFlagSet := flag.NewFlagSet("import", flag.ExitOnError)
	importFormat := importFlagSet.String("format", "csv", "Import content format")

	// Process first argument to program, which should be a sub-command

	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "no sub-command was specified")
		os.Exit(1)
	}

	// Handle sub-commands
	switch os.Args[1] {

	case "connect":
		arg := os.Args[2]
		fmt.Println("connecting to ", arg)

	case "defaults":
	case "list":
	case "def":

	case "a":
		addFlagSet.Parse(os.Args[2:])
		add(addCmdCnSettings)

	case "add":
		addFlagSet.Parse(os.Args[2:])
		add(addCmdCnSettings)

	case "s":
		setFlagSet.Parse(os.Args[2:])
		set(*setID, setCmdCnSettings)

	case "set":
		setFlagSet.Parse(os.Args[2:])
		set(*setID, setCmdCnSettings)

	case "rm":
		arg := os.Args[2]
		fmt.Println("Deleting ", arg)

	case "export":
		exportFlagSet.Parse(os.Args[2:])
		fmt.Println("export connections.")
		fmt.Println("    format:", *exportFormat)

	case "import":
		importFlagSet.Parse(os.Args[2:])
		fmt.Println("Import connections.")
		fmt.Println("    format:", *importFormat)

	case "search":
		arg := os.Args[2]
		fmt.Println("Searching for ", arg)

	case "help":
		fmt.Println(helpBlurb())
		os.Exit(0)

	default:
		fmt.Fprintln(os.Stderr, "unknown sub-command")
		os.Exit(1)
	}

}
