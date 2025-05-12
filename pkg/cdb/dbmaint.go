package cdb

import (
	"database/sql"
	"strings"

	"golang.org/x/mod/semver"
)

const schemaVersion = "v1.1"

var schemas = map[string]string{
	"v1.0": `
		CREATE TABLE 'global' (
			'setting'   TEXT UNIQUE,
			'value'     TEXT,
			PRIMARY KEY('setting')
		);
		CREATE TABLE 'defaults' (
			'setting'       TEXT UNIQUE,
			'value'         TEXT,
			PRIMARY KEY('setting')
		);
		CREATE TABLE 'connections' (
			'id'            INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
			'nickname'      TEXT NOT NULL UNIQUE,
			'host'          TEXT NOT NULL,
			'user'          TEXT,
			'description'   TEXT,
			'args'          TEXT,
			'identity'      TEXT,
			'command'       TEXT
		);`,
	"v1.1": `
		CREATE TABLE 'global' (
			'setting'	TEXT UNIQUE,
			'value'	    TEXT,
			PRIMARY KEY('setting')
		);
		CREATE TABLE 'defaults' (
			'setting'	TEXT UNIQUE,
			'value'	    TEXT,
			PRIMARY KEY('setting')
		);
		CREATE TABLE 'connections' (
			'id'         INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
			'nickname'      TEXT NOT NULL UNIQUE,
			'host'          TEXT NOT NULL,
			'user'          TEXT,
			'description'   TEXT,
			'args'          TEXT,
			'identity'      TEXT,
			'command'       TEXT,
			'binary'        TEXT
		);`,
}

var schemaUpgrades = map[string]string{
	"v1.1": `
		ALTER TABLE 'connections' ADD COLUMN 'binary' TEXT;
		INSERT INTO 'defaults' (setting,value) VALUES ('binary',NULL);
		UPDATE 'global' SET value='1.1' WHERE setting='schema_version';
	);`,
}

// createDb will populate an empty Sqlite file with the tables and default
// values sshcm expects.
// The function will return nil upon completion or an error when an exception
// occurs.
func createDb(db *sql.DB, version string) error {
	if !semver.IsValid(version) {
		return ErrSchemaVerInvalid
	}

	// Check schemaVer
	_, ok := schemas[version]

	if !ok {
		return ErrSchemaVerInvalid
	}

	// Create table schema
	_, err := db.Exec(schemas[version])

	if err != nil {
		return err
	}

	// Initialize global settings
	_, err = db.Exec(`
			INSERT INTO 'global' (setting,value)
			VALUES ('schema_version',$1);
		`, version)

	if err != nil {
		return err
	}

	// Initialize default options
	_, err = db.Exec(`
			BEGIN TRANSACTION;
			INSERT INTO 'defaults' (setting,value) VALUES ('binary',NULL);
			INSERT INTO 'defaults' (setting,value) VALUES ('user',NULL);
			INSERT INTO 'defaults' (setting,value) VALUES ('args',NULL);
			INSERT INTO 'defaults' (setting,value) VALUES ('identity',NULL);
			INSERT INTO 'defaults' (setting,value) VALUES ('command',NULL);
			COMMIT;
		`)

	return err
}

// getDbSchemaVersion will read and return the schema version from an sshcm
// Sqlite database file.
func getDbSchemaVersion(db *sql.DB) (string, error) {
	var v sql.NullString

	row := db.QueryRow(
		`SELECT value
		FROM global
		WHERE setting = 'schema_version'`)
	err := row.Scan(&v)

	// Run some checks against the version number we pulled from the database
	if err != nil {
		// We got an error running the query
		return "", err
	} else if !v.Valid {
		// The read version string is invalid
		return "", ErrSchemaVerInvalid
	} else if !strings.HasPrefix(v.String, "v") {
		// The version number doesn't start with a "v"
		// This might be recoverable, in that it might be a DB from the Tcl version

		v.String = "v" + v.String
	}

	// Final check - is this a valid semantic version?
	if !semver.IsValid(v.String) {
		return "", ErrSchemaVerInvalid
	}

	return v.String, nil
}

// isDbCurrent calls getDbSchemaVersion to read the schema version from an
// sshcm connection database and compares it to the version expected by the
// specific version of this library.
//
// If the database version matches the version expected by this library, the
// function will return true, nil.
//
// If the database version does not match the  version expected by this library
// the function will return false, nil.
//
// If an error occurs, the function will return false, error.
func isDbCurrent(db *sql.DB) (bool, error) {
	dbVer, err := getDbSchemaVersion(db)

	if err != nil {
		return false, err
	}

	if semver.Compare(schemaVersion, dbVer) == 0 {
		return true, nil
	}

	return false, nil
}

// isDbSchemaVersionSupported checks whether the sshcm DB schema version is
// supported by this library. Support, in this context, means that the DB schema
// is either the specific version required by this library or the library can
// perform an upgrade. This nuance can be determined by using this function in
// concert with isDbCurrent.
//
// If the DB schema is supported, this function will return true, otherwise
// it will return false.
//
// If an error occurs, err will be non-nil.
func isDbSchemaVersionSupported(db *sql.DB) (bool, error) {
	// TODO: Rewrite this and isDbCurrent to be more explicit about
	// upgradeability circumstances.
	dbVer, err := getDbSchemaVersion(db)

	if err != nil {
		return false, err
	}

	// Check schemaVer
	_, ok := schemas[dbVer]

	return ok, nil
}

// upgradeDbSchema upgrades an sshcm connection database to the version this
// library supports.
//
// NOTE: This feature does nothing, currently.
func upgradeDbSchema(db *sql.DB) error {
	// TODO: Determine upgrade strategy (ex. do we have to do multiple upgrades?)

	// TODO: Do upgrades

	return nil
}
