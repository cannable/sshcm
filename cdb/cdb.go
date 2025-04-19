package cdb

import (
	"database/sql"
	"errors"
	"os"
	"os/user"
	"strings"

	_ "modernc.org/sqlite"
)

type ConnectionDB struct {
	connection *sql.DB
	Defaults   map[string]string
}

type Connection struct {
	Id          int64
	Nickname    string
	Host        string
	User        string
	Description string
	Args        string
	Identity    string
	Command     string
	Binary      string
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

	// Create tables
	if create {
		if _, err = db.Exec(`
			BEGIN TRANSACTION;
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
			);
			INSERT INTO 'global' (setting,value) VALUES ('schema_version','1.1');
			COMMIT;
		`); err != nil {
			return cdb, err
		}

		// Set default options
		if _, err = db.Exec(`
			BEGIN TRANSACTION;
			INSERT INTO 'defaults' (setting,value) VALUES ('binary',NULL);
			INSERT INTO 'defaults' (setting,value) VALUES ('user',NULL);
			INSERT INTO 'defaults' (setting,value) VALUES ('args',NULL);
			INSERT INTO 'defaults' (setting,value) VALUES ('identity',NULL);
			INSERT INTO 'defaults' (setting,value) VALUES ('command',NULL);
			COMMIT;
		`); err != nil {
			return cdb, err
		}
	}

	// Assemble connection struct & prep for loading default settings
	cdb.connection = db
	defaults := make(map[string]string)

	// Set the default user to the current one before loading application defaults
	u, err := user.Current()

	if err == nil {
		defaults["user"] = u.Username
	}

	// Read default settings from DB
	rows, err := db.Query("SELECT setting,value FROM defaults")

	if err != nil {
		return cdb, err
	}

	for rows.Next() {
		var k, v sql.NullString

		if err := rows.Scan(&k, &v); err != nil {
			return cdb, err
		}

		if v.Valid {
			defaults[k.String] = v.String
		}
	}

	cdb.Defaults = defaults

	return cdb, err
}

func marshallConnection(row *sql.Row, def map[string]string) (Connection, error) {
	var c Connection
	var id sql.NullInt64
	var nickname, host, user, description, args, identity, command, binary sql.NullString

	// Get connection details from DB
	err := row.Scan(&id, &nickname, &host, &user, &description, &args, &identity, &command, &binary)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c, ErrMarshallNoRows
		}
		return c, err
	}

	// Defaults
	c.User = def["user"]
	c.Args = def["args"]
	c.Identity = def["identity"]
	c.Command = def["command"]
	c.Binary = def["binary"]

	// If we're missing an id, nickname, or host name, something's very wrong
	if !(id.Valid || nickname.Valid || host.Valid) {
		return c, err
	}

	c.Id = id.Int64
	c.Nickname = nickname.String
	c.Host = host.String

	// If the DB has a user name for this connection, use it
	if user.Valid {
		c.User = user.String
	}

	// The remaining bits from the DB are optional

	// Use the description from the DB if it exists
	if description.Valid {
		c.Description = description.String
	}

	// Use the args from the DB if it exists
	if args.Valid {
		c.Args = args.String
	}

	// Use the identity from the DB if it exists
	if identity.Valid {
		c.Identity = identity.String
	}

	// Use the command from the DB if it exists
	if command.Valid {
		c.Command = command.String
	}

	// Use the binary from the DB if it exists
	if binary.Valid {
		c.Binary = binary.String
	}

	return c, err
}

func (conndb *ConnectionDB) Close() {
	defer conndb.connection.Close()
}

func (conndb *ConnectionDB) IdExists(id int64) bool {
	var check int

	err := conndb.connection.QueryRow("SELECT id FROM connections WHERE id = $1", id).Scan(&check)

	return err == nil
}

func sqlNullableString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}

	return sql.NullString{String: s, Valid: true}
}

func (conndb *ConnectionDB) NicknameExists(nickname string) bool {
	var check int

	err := conndb.connection.QueryRow("SELECT nickname FROM connections WHERE nickname = $1", nickname).Scan(&check)

	return err == nil
}

func (conndb *ConnectionDB) GetId(nickname string) (int, error) {
	var id int

	err := conndb.connection.QueryRow("SELECT id FROM connections WHERE nickname = $1", nickname).Scan(&id)

	return id, err
}

func (conndb *ConnectionDB) Get(id int64) (Connection, error) {
	var c Connection

	// Get connection details from DB
	row := conndb.connection.QueryRow(`
		SELECT
			id,
			nickname,
			host,
			user,
			description,
			args,
			identity,
			command,
			binary
		FROM connections
		WHERE id = $1
	`, id)

	c, err := marshallConnection(row, conndb.Defaults)

	return c, err
}

func isValidProperty(property string) bool {
	valid := [8]string{
		"nickname",
		"host",
		"user",
		"description",
		"args",
		"identity",
		"command",
		"binary"}

	for _, v := range valid {
		if strings.Compare(property, v) == 0 {
			return true
		}
	}
	return false
}

func (conndb *ConnectionDB) GetByProperty(property string, value string) (Connection, error) {
	var c Connection

	if !isValidProperty(property) {
		return c, ErrPropertyInvalid
	}

	// Get connection details from DB
	row := conndb.connection.QueryRow(`
		SELECT
			id,
			nickname,
			host,
			user,
			description,
			args,
			identity,
			command,
			binary
		FROM connections
		WHERE `+property+" = $1", value)

	c, err := marshallConnection(row, conndb.Defaults)

	return c, err
}

func (conndb *ConnectionDB) Update(c *Connection) error {
	// Do we have an id?
	if c.Id < 1 {
		return ErrUpdateNoId
	}

	// Does the ID exist?
	if !conndb.IdExists(c.Id) {
		return ErrUpdateIdNotExist
	}

	// Do we have a nickname?
	if len(c.Nickname) < 1 {
		return ErrUpdateNoNickname
	}

	// Try updating the connection
	_, err := conndb.connection.Exec(`
		UPDATE connections SET
			nickname = $2,
			host = $3,
			user = $4,
			description = $5,
			args = $6,
			identity = $7,
			command = $8,
			binary = $9
		WHERE id = $1
		`,
		c.Id,
		c.Nickname,
		c.Host,
		sqlNullableString(c.User),
		sqlNullableString(c.Description),
		sqlNullableString(c.Args),
		sqlNullableString(c.Identity),
		sqlNullableString(c.Command),
		sqlNullableString(c.Binary))

	return err
}

func (conndb *ConnectionDB) Add(c *Connection) (int64, error) {
	// Do we have a nickname?
	if len(c.Nickname) < 1 {
		return -1, ErrAddNoNickname
	}

	// See if the nickname already exists
	if conndb.NicknameExists(c.Nickname) {
		return -1, ErrAddNicknameExists
	}

	// Try adding the connection
	result, err := conndb.connection.Exec(`
		INSERT INTO connections (
			nickname,
			host,
			user,
			description,
			args,
			identity,
			command,
			binary
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8
		)`,
		c.Nickname,
		c.Host,
		sqlNullableString(c.User),
		sqlNullableString(c.Description),
		sqlNullableString(c.Args),
		sqlNullableString(c.Identity),
		sqlNullableString(c.Command),
		sqlNullableString(c.Binary))

	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return -1, err
	}

	return id, err
}

func (conndb *ConnectionDB) GetAll() ([]*Connection, error) {
	var cns []*Connection
	rows, err := conndb.connection.Query("SELECT id FROM connections")

	if err != nil {
		return cns, err
	}

	for rows.Next() {
		var id int64

		if err := rows.Scan(&id); err != nil {
			return cns, err
		}

		c, err := conndb.Get(id)

		if err != nil {
			return cns, err
		}

		cns = append(cns, &c)
	}

	return cns, err
}
