package cdb

import (
	"fmt"
)

type Connection struct {
	db          *ConnectionDB
	Id          *IdProperty
	Nickname    *NicknameProperty
	Host        *ConnectionProperty
	User        *ConnectionProperty
	Description *ConnectionProperty
	Args        *ConnectionProperty
	Identity    *ConnectionProperty
	Command     *ConnectionProperty
	Binary      *ConnectionProperty
}

func (c *Connection) Delete() error {
	// See if this connection has a parent ConnectionDB attached
	if c.db == nil {
		return ErrConnNoDb
	}

	// Do we have an id?
	if c.Id.IsEmpty() {
		return ErrConnNoId
	}

	// Does the ID exist?
	if !c.db.Exists(c.Id.Value) {
		return ErrIdNotExist
	}

	// Try deleting the connection
	_, err := c.db.connection.Exec(`
        DELETE FROM connections
		WHERE id = $1
		`,
		c.Id.SqlNullableValue())

	return err
}

func (c *Connection) Update() error {
	// See if this connection has a parent ConnectionDB attached
	if c.db == nil {
		return ErrConnNoDb
	}

	// Do we have an id?
	if c.Id.IsEmpty() {
		return ErrConnNoId
	}

	// Does the ID exist?
	if !c.db.Exists(c.Id.Value) {
		return ErrIdNotExist
	}

	// Do we have a nickname?
	if c.Nickname.IsEmpty() {
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
		c.Id.SqlNullableValue(),
		c.Nickname.SqlNullableValue(),
		c.Host.SqlNullableValue(),
		c.User.SqlNullableValue(),
		c.Description.SqlNullableValue(),
		c.Args.SqlNullableValue(),
		c.Identity.SqlNullableValue(),
		c.Command.SqlNullableValue(),
		c.Binary.SqlNullableValue(),
	)

	return err
}

func (c *Connection) TemplateTrimmer(s string, len int) string {
	f := fmt.Sprintf("%-*s", len, s)

	return f[:len]
}
