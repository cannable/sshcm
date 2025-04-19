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

	// Flag definitions - Add connection
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)

	addArgs := addCmd.String("args", "", "Arguments to pass to SSH command")
	addCommand := addCmd.String("command", "", "SSH command to run")
	addDescription := addCmd.String("description", "", "Short description of the connection")
	addHost := addCmd.String("host", "", "Connection hostname (or IP address)")
	addIdentity := addCmd.String("identity", "", "SSH identity to use for connection (a la '-i')")
	addNickname := addCmd.String("nickname", "", "Nickname for connection")
	addUser := addCmd.String("user", "", "User name for connection")

	// Flag definitions - Set connection property
	setCmd := flag.NewFlagSet("set", flag.ExitOnError)

	setArgs := setCmd.String("args", "", "Arguments to pass to SSH command")
	setCommand := setCmd.String("command", "", "SSH command to run")
	setDescription := setCmd.String("description", "", "Short description of the connection")
	setHost := setCmd.String("host", "", "Connection hostname (or IP setress)")
	setID := setCmd.Int("id", 0, "Connection id")
	setIdentity := setCmd.String("identity", "", "SSH identity to use for connection (a la '-i')")
	setNickname := setCmd.String("nickname", "", "Nickname for connection")
	setUser := setCmd.String("user", "", "User name for connection")

	// Flag definitions - Export connections
	exportCmd := flag.NewFlagSet("export", flag.ExitOnError)

	exportFormat := exportCmd.String("format", "csv", "export content format")

	// Flag definitions - Import connections
	importCmd := flag.NewFlagSet("import", flag.ExitOnError)

	importFormat := importCmd.String("format", "csv", "Import content format")

	// Process first argument to program, which should be a sub-command

	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "no sub-command was specified")
		os.Exit(1)
	}

	switch os.Args[1] {

	case "connect":
		arg := os.Args[2]
		fmt.Println("connecting to ", arg)

	case "defaults":
	case "list":
	case "def":

	case "add":
		addCmd.Parse(os.Args[2:])
		//addCmd.PrintDefaults()
		fmt.Println("Add connection.")
		fmt.Println("    arguments:", *addArgs)
		fmt.Println("    command:", *addCommand)
		fmt.Println("    description:", *addDescription)
		fmt.Println("    host:", *addHost)
		fmt.Println("    identity:", *addIdentity)
		fmt.Println("    nickname:", *addNickname)
		fmt.Println("    user:", *addUser)

		// TODO: Check that connection nickname starts with a letter

	case "set":
		setCmd.Parse(os.Args[2:])
		fmt.Println("Set connection.")
		fmt.Println("    arguments:", *setArgs)
		fmt.Println("    command:", *setCommand)
		fmt.Println("    description:", *setDescription)
		fmt.Println("    host:", *setHost)
		fmt.Println("    id:", *setID)
		fmt.Println("    identity:", *setIdentity)
		fmt.Println("    nickname:", *setNickname)
		fmt.Println("    user:", *setUser)

	case "rm":
		arg := os.Args[2]
		fmt.Println("Deleting ", arg)

	case "export":
		exportCmd.Parse(os.Args[2:])
		fmt.Println("export connections.")
		fmt.Println("    format:", *exportFormat)

	case "import":
		importCmd.Parse(os.Args[2:])
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
