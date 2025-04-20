package cdb

type Connection struct {
	db          *ConnectionDB
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

func (c *Connection) Delete() error {
	// See if this connection has a parent ConnectionDB attached
	if c.db == nil {
		return ErrConnNoDb
	}

	// Do we have an id?
	if c.Id < 1 {
		return ErrConnNoId
	}

	// Does the ID exist?
	if !c.db.Exists(c.Id) {
		return ErrIdNotExist
	}

	// Try deleting the connection
	_, err := c.db.connection.Exec(`
        DELETE FROM connections
		WHERE id = $1
		`,
		c.Id)

	return err
}

func (c *Connection) Update() error {
	// See if this connection has a parent ConnectionDB attached
	if c.db == nil {
		return ErrConnNoDb
	}

	// Do we have an id?
	if c.Id < 1 {
		return ErrConnNoId
	}

	// Does the ID exist?
	if !c.db.Exists(c.Id) {
		return ErrIdNotExist
	}

	// Do we have a nickname?
	if len(c.Nickname) < 1 {
		return ErrConnNoNickname
	}

	// Try updating the connection
	_, err := c.db.connection.Exec(`
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
