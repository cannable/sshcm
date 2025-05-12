package cdb

import (
	"database/sql"
	"strings"

	"golang.org/x/mod/semver"
)

const SchemaVersion = "v1.1"

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

// CheckDbHealth runs health checks on the connection DB and returns an error
// if there is an issue.
func (conndb *ConnectionDB) CheckDbHealth() error {
	// Read the DB version
	version, err := conndb.GetDbSchemaVersion()

	if err != nil {
		return err
	}

	// Return validation checks
	return ValidateDbSchemaVersion(version)
}

// InitializeDb will populate an empty Sqlite file with the tables and default
// values sshcm expects.
// The function will return nil upon completion or an error when an exception
// occurs.
func (conndb *ConnectionDB) InitializeDb(version string) error {
	err := ValidateDbSchemaVersion(version)

	db := conndb.connection

	if err != nil {
		return err
	}

	// Create table schema
	_, err = db.Exec(schemas[version])

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

// GetDbSchemaVersion will read and return the schema version from an sshcm
// Sqlite database file. This does not validate whether the schema version is
// usable by this package, it simply reads the version from the DB and returns
// the result.
func (conndb *ConnectionDB) GetDbSchemaVersion() (string, error) {
	var v sql.NullString

	db := conndb.connection

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
	}

	return v.String, nil
}

// ValidateDbSchemaVersion runs checks against the passed schema version to see
// if it's supported by this package. It will return nil if the DB is usable,
// and various errors otherwise:
//
//		ErrSchemaVerInvalid - Unrecoverable. Something unexpected happened.
//		ErrSchemaTooOld - Recoverable. The caller should call upgradeDbSchema() to
//			attempt to upgrade the DB schema to the latest version.
//		ErrSchemaTooNew - Unrecoverable. This package (or the calling tool) needs
//	   to be upgraded.
//	 Others - Likely unrecoverable. Other errors returned by called funcs.
func ValidateDbSchemaVersion(version string) error {
	if !strings.HasPrefix(version, "v") {
		// The version number doesn't start with a "v"
		// This might be recoverable, in that it might be a DB from the Tcl version

		version = "v" + version
	}

	// Is this a valid semantic version?
	if !semver.IsValid(version) {
		return ErrSchemaVerInvalid
	}

	// Compare the version number to this tool
	compare := semver.Compare(version, SchemaVersion)

	if compare == 0 {
		// If the DB is the same version as this tool supports, we're done with the
		// checks and good to go!
		return nil
	} else if compare > 0 {
		// The DB is newer than this tool
		return ErrSchemaTooNew
	}

	// The DB schema version is too old. See if an upgrade is supported.
	// Check schemaVer
	_, ok := schemas[version]

	if ok {
		return ErrSchemaUpgradeNeeded
	}

	// This tool can't upgrade the schema
	return ErrSchemaNoUpgrade
}
