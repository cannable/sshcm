package cdb

import (
	"database/sql"
	"errors"
	"strconv"
)

type ConnectionDB struct {
	connection DbConnIface
}

// DbConnIface provides an interface for interacting with a DB (or mock)
type DbConnIface interface {
	Close() error
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

func (conndb *ConnectionDB) Add(c *Connection) (int64, error) {
	err := c.Validate()

	// The only error we should get from validation is that the connection ID is zero.
	if err != ErrConnIdZero {
		return -1, err
	}

	// See if the nickname already exists
	exists, err := conndb.ExistsByProperty("nickname", c.Nickname)

	if err != nil {
		return -1, err
	} else if exists {
		return -1, ErrDuplicateNickname
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
			command
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7
		)`,
		sqlNullableString(c.Nickname),
		sqlNullableString(c.Host),
		sqlNullableString(c.User),
		sqlNullableString(c.Description),
		sqlNullableString(c.Args),
		sqlNullableString(c.Identity),
		sqlNullableString(c.Command),
	)

	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return -1, err
	}

	return id, err
}

func (conndb *ConnectionDB) Exists(id int64) (bool, error) {
	var check int

	err := conndb.connection.QueryRow("SELECT id FROM connections WHERE id = $1", id).Scan(&check)

	// Return true/false explicitly (based on specific success/fail conditions)
	if err == nil {
		return true, nil
	} else if err == sql.ErrNoRows {
		return false, nil
	}

	return false, err
}

func (conndb *ConnectionDB) ExistsByProperty(property string, value string) (bool, error) {
	// Make sure we're dealing with a valid property first
	if !IsValidProperty(property) {
		return false, ErrInvalidConnectionProperty
	}

	var check int

	err := conndb.connection.QueryRow(`
		SELECT id
		FROM connections
		WHERE `+property+" = $1", value).Scan(&check)

	// Return true/false explicitly (based on specific success/fail conditions)
	if err == nil {
		return true, nil
	} else if err == sql.ErrNoRows {
		return false, nil
	}

	return false, err
}

func (conndb *ConnectionDB) Get(id int64) (Connection, error) {
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
			command
		FROM connections
		WHERE id = $1
		ORDER BY id;
	`, id)

	// Get connection details from DB
	var sqlId sql.NullInt64
	var nickname, host, user, description, args, identity, command sql.NullString

	err := row.Scan(
		&sqlId,
		&nickname,
		&host,
		&user,
		&description,
		&args,
		&identity,
		&command,
	)

	// Check SQL scanning errors before continuing
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Connection{}, ErrConnectionNotFound
		}

		return Connection{}, err
	}

	// If any of the connection properties are invalid, fail
	if !(sqlId.Valid ||
		nickname.Valid ||
		host.Valid ||
		user.Valid ||
		description.Valid ||
		args.Valid ||
		identity.Valid ||
		command.Valid) {
		return Connection{}, ErrConnFromDbInvalid
	}

	// Attach the connection to its parent (so that connection methods work)
	c := Connection{
		db:          conndb,
		Id:          sqlId.Int64,
		Nickname:    nickname.String,
		Host:        host.String,
		User:        user.String,
		Description: description.String,
		Args:        args.String,
		Identity:    identity.String,
		Command:     command.String,
	}

	err = c.Validate()

	return c, err
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

// GetByIdOrNickname looks up a connection by id or nickname, then returns a
// Connection struct.
// If the look up succeeded, err will be nil and it can be assumed that the
// Connection is safe to use.
func (conndb *ConnectionDB) GetByIdOrNickname(arg string) (Connection, error) {
	var c Connection

	// Get connection by ID or nickname
	if err := ValidateId(arg); err == nil {
		// Got a valid ID
		id, err := strconv.Atoi(arg)

		if err != nil {
			return c, err
		}

		// Get connection by id
		c, err = conndb.Get(int64(id))

		if err != nil {
			return c, err
		}
	} else {
		err := ValidateNickname(arg)

		if err == nil {
			// Got a valid nickname
			nickname := arg

			// Get connection by nickname
			c, err = conndb.GetByProperty("nickname", nickname)

			if err != nil {
				return c, err
			}
		} else {
			return c, ErrInvalidIdOrNickname
		}
	}

	return c, nil
}

func (conndb *ConnectionDB) GetByProperty(property string, value string) (Connection, error) {
	if !IsValidProperty(property) {
		return Connection{}, ErrPropertyInvalid
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
			command
		FROM connections
		WHERE `+property+" = $1", value)

	// Get connection details from DB
	var sqlId sql.NullInt64
	var nickname, host, user, description, args, identity, command sql.NullString

	err := row.Scan(
		&sqlId,
		&nickname,
		&host,
		&user,
		&description,
		&args,
		&identity,
		&command,
	)

	// Check SQL scanning errors before continuing
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Connection{}, ErrConnectionNotFound
		}

		return Connection{}, err
	}

	// If any of the connection properties are invalid, fail
	if !(sqlId.Valid ||
		nickname.Valid ||
		host.Valid ||
		user.Valid ||
		description.Valid ||
		args.Valid ||
		identity.Valid ||
		command.Valid) {
		return Connection{}, ErrConnFromDbInvalid
	}

	// Attach the connection to its parent (so that connection methods work)
	c := Connection{
		db:          conndb,
		Id:          sqlId.Int64,
		Nickname:    nickname.String,
		Host:        host.String,
		User:        user.String,
		Description: description.String,
		Args:        args.String,
		Identity:    identity.String,
		Command:     command.String,
	}

	err = c.Validate()

	return c, err
}

func (conndb *ConnectionDB) Search(search string) ([]*Connection, error) {
	var cns []*Connection
	rows, err := conndb.connection.Query(`
		SELECT id
		FROM connections
		WHERE (nickname LIKE $1)
		OR (host LIKE $1)
		OR (user LIKE $1)
		OR (description LIKE $1)
		ORDER BY id;
	`, "%"+search+"%")

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
