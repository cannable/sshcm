package cdb

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

func (conndb *ConnectionDB) Close() {
	defer conndb.connection.Close()
}

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
		// Can we use the DB?
		dbUsable, err := isDbSchemaVersionSupported(db)

		if err != nil {
			return cdb, err
		}

		if !dbUsable {
			return cdb, ErrDbVersionNotRecognized
		}

		// See if the DB schema is old
		dbLatest, err := isDbCurrent(db)

		if err != nil {
			return cdb, err
		}

		// Do schema upgrade, if needed
		if !dbLatest {
			err = upgradeDbSchema(db)

			if err != nil {
				return cdb, err
			}
		}
	}

	// Assemble connection struct & prep for loading default settings
	cdb.connection = db

	return cdb, nil
}
