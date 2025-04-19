package cdb

import "database/sql"

type ConnectionDB struct {
	connection *sql.DB
	Defaults   map[string]string
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

func (conndb *ConnectionDB) Exists(id int64) bool {
	var check int

	err := conndb.connection.QueryRow("SELECT id FROM connections WHERE id = $1", id).Scan(&check)

	return err == nil
}

func (conndb *ConnectionDB) ExistsByProperty(property string, value string) bool {

	if !isValidProperty(property) {
		return false
	}

	var check int

	err := conndb.connection.QueryRow(`
		SELECT COUNT(nickname)
		FROM connections
		WHERE`+property+" = $1", value).Scan(&check)

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
			command,
			binary
		FROM connections
		WHERE id = $1
	`, id)

	c, err := marshallConnection(row, conndb.Defaults)

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

	// Attach the connection to its parent (so that connection methods work)
	c.db = conndb

	return c, err
}
