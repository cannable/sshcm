package cdb

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
)

// A Connection is an SSH connection, as stored in a ConnectionDB.
// A Connection could originate from a ConnectionDB (ex. returned by a SELECT
// from the underlying SQL DB), or could be destined for one (ex. for later
// INSERTion).
//
// If a connection originated from a ConnectionDB, a pointer to it will be
// stored in db. This allows methods such as Update to perform database
// operations from the context of the Connection. Connections originating
// from a ConnectionDB will also have an Id set.
//
// Connections created detached from a ConnectionDB will not have a db or Id
// set. This is a primitive safety mechanism to prevent arbitrary writes to the
// database bypassing validations that avoid throwing SQL errors (like
// checking for Nickname uniqueness).
type Connection struct {
	db          *ConnectionDB       // pointer to parent ConnectionDB
	Id          *IdProperty         // unique connection id
	Nickname    *NicknameProperty   // unique connection nickname
	Host        *ConnectionProperty // connection-specific host name/IP address
	User        *ConnectionProperty // connection-specific user name
	Description *ConnectionProperty // connection-specific description (hopefully friendly)
	Args        *ConnectionProperty // connection-specific arguments to pass to SSH Command
	Identity    *ConnectionProperty // connection-specific OpenSSH-style identity string (ex. path or name)
	Command     *ConnectionProperty // connection-specific Command to run (ex. sftp)
	Binary      *ConnectionProperty // to be deleted
}

// Delete removes a connection from the underlying SQL database.
// It will return nil if the operation succeeeded and err otherwise.
// Several checks are implemented that return package-specific errors. These
// checks are simple and only cover obvious situations that will cause SQL
// query exceptions.
func (c Connection) Delete() error {
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

func (c Connection) WriteCSV(w *csv.Writer) error {
	return w.Write([]string{
		c.Id.String(),
		c.Nickname.String(),
		c.User.String(),
		c.Host.String(),
		c.Description.String(),
		c.Args.String(),
		c.Identity.String(),
		c.Command.String(),
	})
}

func (c Connection) WriteJSON(w io.Writer) error {
	j, err := json.Marshal(c)

	if err != nil {
		return err
	}

	_, err = w.Write(j)

	return err
}

// String returns a string containing the connection nickname and ID, in the
// format: "nickname (id)".
func (c Connection) String() string {
	return fmt.Sprintf("%s (%d)", c.Nickname.Value, c.Id.Value)
}

// Update updates an existing connection in the SQL database.
// It will return nil if the operation succeeeded and err otherwise.
// Several checks are implemented that return package-specific errors. These
// checks are simple and only cover obvious situations that will cause SQL
// query exceptions.
func (c Connection) Update() error {
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
			command = $8
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
	)

	return err
}
