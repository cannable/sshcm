// Package cdb provides a database abstraction layer for storing and retrieving
// SSH connection information for the sshcm utility.
package cdb

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

// Close Gracefully closes a connection to a database.
func (conndb ConnectionDB) Close() {
	defer conndb.connection.Close()
}

// Connect connects to a database and returns a ConnectionDB. In the event an
// error occurs, it'll be returned.
func Connect(driver string, path string) (ConnectionDB, error) {
	var cdb ConnectionDB

	supported := false

	switch driver {

	case "sqlite":
		supported = true
	}

	if !supported {
		return cdb, ErrUnsupportedSqlDriver
	}

	db, err := sql.Open(driver, path)

	if err != nil {
		return cdb, err
	}

	// Assemble connection struct & prep for loading default settings
	cdb.connection = db

	return cdb, nil
}

// NewConnection will create a new, empty Connection struct.
// Using this function is preferred vs. creating a new struct via literal, as
// the format of the struct may change in the future.
func NewConnection() Connection {
	return Connection{}
}
