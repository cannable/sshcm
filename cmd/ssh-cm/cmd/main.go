package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
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
		cmd.PersistentFlags().Int64VarP(&cmdCnId, "id", "i", -1, "ID of connection")
	}
}

func addConnection() (int64, error) {
	// Validate nickname follows the correct convention
	err := cdb.ValidateNickname(setNewNickname)

	if err != nil {
		return -1, err
	}

	// Nicknames must be unique. See if this one exists.
	exists := db.ExistsByProperty("nickname", cmdCnNickname)

	if exists {
		return -1, ErrNicknameExists
	}

	c := cdb.NewConnection()

	c.Nickname.Value = cmdCnNickname
	c.Host.Value = cmdCnHost
	c.User.Value = cmdCnUser
	c.Description.Value = cmdCnDescription
	c.Args.Value = cmdCnArgs
	c.Identity.Value = cmdCnIdentity
	c.Command.Value = cmdCnCommand

	if debugMode {
		fmt.Println("Adding connection:")
		printConnection(&c, false)
	}

	id, err := db.Add(&c)

	if err != nil {
		return -1, err
	}

	return id, nil
}

func connect(arg string) error {
	c, err := getCnByIdOrNickname(arg)

	if err != nil {
		return err
	}

	if debugMode {
		//fmt.Printf("Connecting to %s (%d)...\n", c.Nickname.Value, c.Id.Value)
		fmt.Println("Connecting to ", c)
	}

	// Get effective SSH command (binary)
	cmd, err := db.GetEffectiveValue(c.Binary.Value, "binary")
	fmt.Println("asdf")

	if err != nil {
		return err
	}

	// If the program default is empty, use 'ssh'
	if strings.Compare(cmd, "") == 0 {
		cmd = "ssh"
	}

	// Make sure ssh binary resolves in PATH
	execBin, err := exec.LookPath("ssh")

	if err != nil {
		return err
	}

	var execArgs = []string{execBin}

	// Append arguments
	args, err := db.GetEffectiveValue(c.Args.Value, "args")

	if err != nil {
		return err
	}

	if strings.Compare(args, "") != 0 {
		// TODO: This is probably really mangled and won't work.
		// Figure out a way to reconstitute flat arguments from the DB.
		execArgs = append(execArgs, args)
	}

	// Append identity
	identity, err := db.GetEffectiveValue(c.Identity.Value, "identity")

	if err != nil {
		return err
	}

	if strings.Compare(identity, "") != 0 {
		execArgs = append(execArgs, "-i", identity)
	}

	// Host & user
	host := c.Host.Value
	user, err := db.GetEffectiveValue(c.User.Value, "user")

	if err != nil {
		return err
	}

	if strings.Compare(user, "") != 0 {
		execArgs = append(execArgs, user+"@"+host)
	} else {
		execArgs = append(execArgs, host)
	}

	if debugMode {
		fmt.Println("connection details:")
		fmt.Printf("binary:   '%s'\n", cmd)
		fmt.Printf("arguments:'%s'\n", execArgs)
	}

	// Connect
	execEnv := os.Environ()

	err = syscall.Exec(execBin, execArgs, execEnv)

	if err != nil {
		return err
	}

	return nil
}

func deleteConnection(arg string) error {
	c, err := getCnByIdOrNickname(arg)

	if err != nil {
		return err
	}

	if debugMode {
		fmt.Printf("Deleting connection %s (%d).\n", c.Nickname, c.Id)
	}

	// Delete connection
	err = c.Delete()

	if err != nil {
		return err
	}

	return nil
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
			return c, cdb.ErrConnNoId
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
			return c, cdb.ErrConnNoNickname
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

func isValidIdOrNickname(s string) bool {
	// Determine if the passed sument is a nickname or id
	if err := cdb.ValidateId(s); err == nil {
		// Got a valid id
		return true
	} else if err := cdb.ValidateNickname(s); err == nil {
		// Got a valid nickname
		return true
	}

	return false
}

func listConnections(cns []*cdb.Connection, wide bool) {
	// Assemble output template
	t := `{{ printf "%-4s" "ID" }} | `
	t = t + `{{ printf "%-15s" "Nickname" }} | `
	t = t + `{{ printf "%-10s" "User" }} | `
	t = t + `{{ printf "%-15s" "Host" }} | `
	t = t + `{{ printf "%-20s" "Description" }} | `

	if wide {
		t = t + `{{ printf "%-10s" "Args" }} | `
		t = t + `{{ printf "%-10s" "Identity" }} | `
		t = t + `{{ printf "%-10s" "Command" }} | `
		t = t + `{{ printf "%-10s"  "Binary" }} | `
	}

	t = t + "\n{{ range . }}"
	t = t + `{{ .Id.StringTrimmed 4 }} | `
	t = t + `{{ .Nickname.StringTrimmed 15 }} | `
	t = t + `{{ .User.StringTrimmed 10 }} | `
	t = t + `{{ .Host.StringTrimmed 15 }} | `
	t = t + `{{ .Description.StringTrimmed 20 }} | `

	if wide {
		t = t + `{{ .Args.StringTrimmed 10 }} | `
		t = t + `{{ .Identity.StringTrimmed 10 }} | `
		t = t + `{{ .Command.StringTrimmed 10 }} | `
		t = t + `{{ .Binary.StringTrimmed 10 }} | `
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

func setConnection() error {
	var c cdb.Connection
	var err error

	// Did we get an ID or nickname?
	if cmdCnId > 0 {
		// Got an ID. Get the connection.
		c, err = db.Get(cmdCnId)

		if err != nil {
			return err
		}
	} else if strings.Compare(cmdCnNickname, "") != 0 {
		// Got a nickname. Get the connection.
		c, err = db.GetByProperty("nickname", cmdCnNickname)

		if err != nil {
			return err
		}
	} else {
		// Got neither... oops
		return ErrNoIdOrNickname
	}

	// Show original values if in debug mode
	if debugMode {
		fmt.Println("Current connection settings:")
		printConnection(&c, false)
	}

	// Determine if we're renaming
	if strings.Compare(setNewNickname, "") != 0 {
		// Validate nickname follows the correct convention
		err = cdb.ValidateNickname(setNewNickname)

		if err != nil {
			return err
		}

		// See if the new nickname exists already.
		exists := db.ExistsByProperty("nickname", setNewNickname)

		if exists {
			return ErrNicknameExists
		}

		c.Nickname = &cdb.NicknameProperty{Value: setNewNickname}
	}

	// Update hostname, if it was passed
	if strings.Compare(cmdCnHost, "") != 0 {
		c.Host.Value = cmdCnHost
	}

	// Update host, if it was passed
	if strings.Compare(cmdCnHost, "") != 0 {
		c.Host.Value = cmdCnHost
	}

	// Update user, if it was passed
	if strings.Compare(cmdCnUser, "") != 0 {
		c.User.Value = cmdCnUser
	}

	// Update description, if it was passed
	if strings.Compare(cmdCnDescription, "") != 0 {
		c.Description.Value = cmdCnDescription
	}

	// Update args, if it was passed
	if strings.Compare(cmdCnArgs, "") != 0 {
		c.Args.Value = cmdCnArgs
	}

	// Update identity, if it was passed
	if strings.Compare(cmdCnIdentity, "") != 0 {
		c.Identity.Value = cmdCnIdentity
	}

	// Update command, if it was passed
	if strings.Compare(cmdCnCommand, "") != 0 {
		c.Command.Value = cmdCnCommand
	}

	// Show to-be-updated values if in debug mode
	if debugMode {
		fmt.Println("\nNew connection settings:")
		printConnection(&c, false)
		fmt.Println("")
	}

	err = c.Update()

	if err != nil {
		return err
	}

	return nil
}

func setDefault(setting string, value string) error {
	return db.SetDefault(setting, value)
}
