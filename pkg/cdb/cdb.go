// Package cdb provides a database abstraction layer for storing and retrieving
// SSH connection information for the sshcm utility.
package cdb

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

// Close Gracefully closes a connection to a database.
func (conndb ConnectionDB) Close() {
	defer conndb.connection.Close()
}

// Open opens a connection to a Sqlite database.
// It returns a new ConnectionDB struct and error.
//
// err will indicate whether a program is able to continue with the database
// connection or not. It does not hint at actions taken to make the database
// usable if it wasn't initially (ex. file creation or schema upgrade).
//
// If the open is successful, err will == nil.
//
// If the database file does not exist, this function will create a new empty
// one one. It will then install the latest table schema and bootstrap
// default setting values. err will still == nil in this case.
//
// If the database file contains an older table schema, this func will upgrade
// it. If the schema is upgraded successfully, err wil also == nil.
func Open(path string) (ConnectionDB, error) {
	var cdb ConnectionDB

	create := false

	// See if calling Open will create a new DB file
	if _, err := os.Stat(path); err != nil {
		create = true
	}

	db, err := sql.Open("sqlite", path)

	if err != nil {
		return cdb, err
	}

	// Create tables, if we need to. If not, see if upgrade is needed
	if create {
		err = createDb(db, schemaVersion)

		if err != nil {
			return cdb, err
		}
	} else {
		// Read the DB version
		version, err := getDbSchemaVersion(db)
		if err != nil {
			return cdb, err
		}

		// Can we use the DB?
		err = validateDbSchemaVersion(version)

		if err != nil {
			if err == ErrSchemaUpgradeNeeded {
				// Do schema upgrade, if needed
				err = upgradeDbSchema(db)

				if err != nil {
					return cdb, err
				}
			} else {
				return cdb, err
			}
		}
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
