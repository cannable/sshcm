package cdb

import (
	"database/sql"
	"strconv"
)

type ConnectionDB struct {
	connection *sql.DB
}

func (conndb *ConnectionDB) Add(c *Connection) (int64, error) {
	// Do we have a nickname?
	if len(c.Nickname) < 1 {
		return -1, ErrConnNoNickname
	}

	// See if the nickname already exists
	if conndb.ExistsByProperty("nickname", c.Nickname) {
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

func (conndb *ConnectionDB) Exists(id int64) bool {
	var check int

	err := conndb.connection.QueryRow("SELECT id FROM connections WHERE id = $1", id).Scan(&check)

	return err == nil
}

func (conndb *ConnectionDB) ExistsByProperty(property string, value string) bool {

	if !IsValidProperty(property) {
		return false
	}

	var check int

	err := conndb.connection.QueryRow(`
		SELECT id
		FROM connections
		WHERE `+property+" = $1", value).Scan(&check)

	return err == nil
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
			command
		FROM connections
		WHERE id = $1
		ORDER BY id;
	`, id)

	c, err := rowToConnection(row)

	// Attach the connection to its parent (so that connection methods work)
	c.db = conndb

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
	var c Connection

	if !IsValidProperty(property) {
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
			command
		FROM connections
		WHERE `+property+" = $1", value)

	c, err := rowToConnection(row)

	// Attach the connection to its parent (so that connection methods work)
	c.db = conndb

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
