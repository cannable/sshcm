package cdb

import (
	"database/sql"
	"strconv"
)

const schemaVersion = 1.1

var schemas = map[float32]string{
	1.0: `
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
	1.1: `
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

var schemaUpgrades = map[float32]string{
	1.1: `
		ALTER TABLE 'connections' ADD COLUMN 'binary' TEXT;
		INSERT INTO 'defaults' (setting,value) VALUES ('binary',NULL);
		UPDATE 'global' SET value='1.1' WHERE setting='schema_version';
	);`,
}

func createDb(db *sql.DB, version float32) error {
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

func getDbSchemaVersion(db *sql.DB) (float32, error) {
	var dbVer sql.NullString

	row := db.QueryRow(
		`SELECT value
		FROM global
		WHERE setting = 'schema_version'`)
	err := row.Scan(&dbVer)

	if err != nil {
		return 0, err
	}

	if !dbVer.Valid {
		return 0, err
	}

	v, err := strconv.ParseFloat(dbVer.String, 32)

	if !dbVer.Valid {
		return 0, err
	}

	return float32(v), err
}

func isDbCurrent(db *sql.DB) (bool, error) {
	dbVer, err := getDbSchemaVersion(db)

	if err != nil {
		return false, err
	}

	if schemaVersion == dbVer {
		return true, nil
	}

	return false, nil
}

func isDbSchemaVersionSupported(db *sql.DB) (bool, error) {
	dbVer, err := getDbSchemaVersion(db)

	if err != nil {
		return false, err
	}

	// Check schemaVer
	_, ok := schemas[dbVer]

	return ok, nil
}

func upgradeDbSchema(db *sql.DB) error {
	// TODO: Determine upgrade strategy (ex. do we have to do multiple upgrades?)

	// TODO: Do upgrades

	return nil
}
