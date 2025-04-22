package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strconv"

	"github.com/cannable/ssh-cm-go/pkg/cdb"
)

func trimmer(s string, len int) string {
	f := fmt.Sprintf("%s%s%s", "%-", strconv.Itoa(len), "s")

	t := fmt.Sprintf(f, s)

	return t[:len]
}

func main() {
	testFilePath := "./cdb-integ-tests.sqlite"

	// Define test connections
	var cnsToAdd = []cdb.Connection{
		cdb.Connection{
			Nickname: "testA",
			Host:     "some.host.name",
		},
		cdb.Connection{
			Nickname: "test_b",
			Host:     "127.0.0.1",
			User:     "tOor",
		},
		cdb.Connection{
			Nickname:    "tes^ C",
			Host:        "blarg",
			Description: "Something profound here",
			User:        "nobody",
		},
		cdb.Connection{
			Nickname: "TESTD",
			Host:     "ggggg",
			Identity: "~/.ssh/id_rsa_demo",
		},
		cdb.Connection{
			Nickname: "teste",
			Host:     "soisoisoi",
			Binary:   "sftp",
		},
	}

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
	db, err := cdb.Open(testFilePath)

	if err != nil {
		log.Fatal("open fail", err)
	}

	// Add test connections
	for _, c := range cnsToAdd {
		id, err := db.Add(&c)

		if err != nil {
			log.Fatal("add: failed ", err)
		}

		log.Println("add: added connection", id)
	}

	// Update a connection
	cn, err := db.Get(2)

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

	fmt.Print(trimmer("Id", 4) + " ")
	fmt.Print(trimmer("Nickname", 15) + " ")
	fmt.Print(trimmer("Host", 20) + " ")
	fmt.Print(trimmer("User", 10) + " ")
	fmt.Print(trimmer("Description", 20) + " ")
	fmt.Print(trimmer("Args", 10) + " ")
	fmt.Print(trimmer("Identity", 10) + " ")
	fmt.Print(trimmer("Command", 10) + " ")
	fmt.Print(trimmer("Binary", 10) + " ")
	fmt.Print("\n")

	for i := range cns {
		idString := strconv.Itoa(int(cns[i].Id))

		fmt.Print(trimmer(idString, 4) + " ")
		fmt.Print(trimmer(cns[i].Nickname, 15) + " ")
		fmt.Print(trimmer(cns[i].Host, 20) + " ")
		fmt.Print(trimmer(cns[i].User, 10) + " ")
		fmt.Print(trimmer(cns[i].Description, 20) + " ")
		fmt.Print(trimmer(cns[i].Args, 10) + " ")
		fmt.Print(trimmer(cns[i].Identity, 10) + " ")
		fmt.Print(trimmer(cns[i].Command, 10) + " ")
		fmt.Print(trimmer(cns[i].Binary, 10) + " ")
		fmt.Print("\n")
	}

	db.Close()

	//Cleanup
	err = os.Remove(testFilePath)

	if err != nil {
		log.Fatal("cleanup: delete failed", err)
	}
	log.Println("cleanup: deleted", testFilePath)
}
