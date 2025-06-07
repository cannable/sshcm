package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/cannable/sshcm/pkg/cdb"
	"github.com/cannable/sshcm/pkg/misc"
)

func trimmer(s string, len int) string {
	f := fmt.Sprintf("%-*s", len, s)

	return f[:len]
}

func main() {
	testFilePath := "./cdb-integ-tests.sqlite"

	// Define test connections
	var cnsToAdd []cdb.Connection

	// Test A
	tst := cdb.NewConnection()
	tst.Nickname = "testA"
	tst.Host = "some.host.name"

	cnsToAdd = append(cnsToAdd, tst)

	// Test B
	tst = cdb.NewConnection()
	tst.Nickname = "test_b"
	tst.Host = "127.0.0.1"
	tst.User = "tOor"

	cnsToAdd = append(cnsToAdd, tst)

	// Test C
	tst = cdb.NewConnection()
	tst.Nickname = "tes^ C"
	tst.Host = "blarg"
	tst.Description = "Something profound here"
	tst.User = "nobody"

	cnsToAdd = append(cnsToAdd, tst)

	// Test D
	tst = cdb.NewConnection()
	tst.Nickname = "TESTD"
	tst.Host = "ggggg"
	tst.Identity = "~/.ssh/id_rsa_demo"

	cnsToAdd = append(cnsToAdd, tst)

	// Test E
	tst = cdb.NewConnection()
	tst.Nickname = "teste"
	tst.Host = "soisoisoi"
	tst.User = "sftp"

	cnsToAdd = append(cnsToAdd, tst)

	// Begin tests
	_, err := os.Stat(testFilePath)

	// File exists, delete it
	if err == nil {
		log.Println("stat: found existing file", testFilePath)

		err = os.Remove(testFilePath)

		if err != nil {
			log.Fatal("stat: delete failed", err)
		}
		log.Println("stat: deleted", testFilePath)
	}

	// Bail on some other kind of error
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			log.Fatal("stat fail", err)
		}
	}

	// Open DB (this should create a new one)
	db, err := cdb.Connect("sqlite", testFilePath)

	if err != nil {
		panic(err)
	}

	err = db.InitializeDb(cdb.SchemaVersion)

	if err != nil {
		log.Fatal("open fail", err)
	}

	// Set & get defaults
	err = db.SetDefault("user", "asdf")

	if err != nil {
		log.Fatal("default: failed to set default user.", err)
	}

	user, err := db.GetDefault("user")

	if err != nil {
		log.Fatal("default: failed to get default user.", err)
	}

	log.Printf("default: set user to '%s'.\n", user)

	// Add test connections
	for _, c := range cnsToAdd {
		id, err := db.Add(&c)

		if err != nil {
			log.Fatal("add: failed ", err)
		}

		log.Println("add: added connection", id)
	}

	log.Println("blarg")

	// Update a connection
	cn, err := db.Get(2)
	log.Println("post-blarg")

	if err != nil {
		log.Fatal("update: failed to get connection ", err)
	}

	cn.Description = "Updated"
	cn.Host = "asdf"
	err = cn.Update()

	if err != nil {
		log.Fatal("update: failed to update connection ", err)
	}

	log.Println("update: updated connection", cn.Id)

	// Test getting a connection by id
	if db.Exists(1) {
		log.Println("exists: got id", 1)
	} else {
		log.Fatal("exists: failed to get id ", 1, err)
	}

	// Test getting a non-existent connection by id
	if db.Exists(99) {
		log.Fatal("exists: got non-existent id ", 99, err)
	} else {
		log.Println("exists: didn't get non-existent id", 99)
	}

	// Test getting id by nickname
	cnByNick, err := db.GetByProperty("nickname", "teste")

	if err != nil {
		log.Fatal("get by nickname: failed ", err)
	} else {
		log.Println("get by nickname: found", cnByNick.Id)
	}

	// Test getting id by non-existent nickname
	cnByNick, err = db.GetByProperty("nickname", "does not exist")

	if err != nil {
		if errors.Is(err, cdb.ErrConnectionNotFound) {
			log.Println("get by nickname: didn't get non-existent nickname", cnByNick.Id)
		} else {
			log.Fatal("get by nickname: ", err)
		}

	} else {
		log.Fatal("get by nickname: found non-existent nickname ", err)
	}

	// Delete a connection
	delcon, err := db.GetByProperty("nickname", "TESTD")

	if err != nil {
		log.Fatal("delete: failed to find test d", err)
	}

	err = delcon.Delete()

	if err != nil {
		log.Fatal("delete: failed to delete test d ", err)
	} else {
		log.Println("delete: deleted test d")
	}

	// Get and print all connections
	cns, err := db.GetAll()

	if err != nil {
		log.Fatal("get all failed ", err)
		os.Exit(1)
	}

	_, err = fmt.Fprintf(os.Stdout, "%s %s %s %s %s %s %s %s\n",
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
		log.Fatal("Error:", err)
		os.Exit(1)
	}

	for i := range cns {
		cns[i].WriteLineLong(os.Stdout)
	}

	db.Close()

	//Cleanup
	err = os.Remove(testFilePath)

	if err != nil {
		log.Fatal("cleanup: delete failed", err)
	}
	log.Println("cleanup: deleted", testFilePath)
}
