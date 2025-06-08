package cdb

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/cannable/sshcm/pkg/misc"
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
	db          *ConnectionDB // pointer to parent ConnectionDB
	Id          int64         // unique connection id
	Nickname    string        // unique connection nickname
	Host        string        // connection-specific host name/IP address
	User        string        // connection-specific user name
	Description string        // connection-specific description (hopefully friendly)
	Args        string        // connection-specific arguments to pass to SSH Command
	Identity    string        // connection-specific OpenSSH-style identity string (ex. path or name)
	Command     string        // connection-specific Command to run (ex. sftp)
	Binary      string        // to be deleted
}

var ListViewColumnWidths = map[string]int{
	"id":          4,
	"nickname":    15,
	"user":        10,
	"host":        15,
	"description": 20,
	"args":        10,
	"identity":    10,
	"command":     10,
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
	if c.Id == 0 {
		return ErrConnNoId
	}

	// Does the ID exist?
	exists, err := c.db.Exists(c.Id)

	if err != nil {
		return err
	}

	if !exists {
		return ErrIdNotExist
	}

	// Try deleting the connection
	_, err = c.db.connection.Exec(`
        DELETE FROM connections
		WHERE id = $1
		`,
		sqlNullableInt64(c.Id))

	return err
}

// WriteRecordLong writes a record-format, multi-line string to the passed
// writer interface. This func will write all connection properties.
// An error will be returned if one occurs, otherwise error will be nil.
func (c Connection) WriteRecordLong(w io.Writer) error {
	offset := 12

	var b strings.Builder

	fmt.Fprintf(&b, "%-*s: %d\n", offset, "ID", c.Id)
	fmt.Fprintf(&b, "%-*s: %s\n", offset, "Nickname", c.Nickname)
	fmt.Fprintf(&b, "%-*s: %s\n", offset, "User", c.User)
	fmt.Fprintf(&b, "%-*s: %s\n", offset, "Host", c.Host)
	fmt.Fprintf(&b, "%-*s: %s\n", offset, "Description", c.Description)
	fmt.Fprintf(&b, "%-*s: %s\n", offset, "Args", c.Args)
	fmt.Fprintf(&b, "%-*s: %s\n", offset, "Identity", c.Identity)
	fmt.Fprintf(&b, "%-*s: %s\n", offset, "Command", c.Command)

	_, err := fmt.Fprint(w, b.String())

	return err
}

// WriteRecordShort writes a record-format, multi-line string to the passed
// writer interface. This func will write only some connection properties.
// An error will be returned if one occurs, otherwise error will be nil.
func (c Connection) WriteRecordShort(w io.Writer) error {
	offset := 12

	var b strings.Builder

	fmt.Fprintf(&b, "%-*s: %d\n", offset, "ID", c.Id)
	fmt.Fprintf(&b, "%-*s: %s\n", offset, "Nickname", c.Nickname)
	fmt.Fprintf(&b, "%-*s: %s\n", offset, "User", c.User)
	fmt.Fprintf(&b, "%-*s: %s\n", offset, "Host", c.Host)
	fmt.Fprintf(&b, "%-*s: %s\n", offset, "Description", c.Description)

	_, err := fmt.Fprint(w, b.String())

	return err
}

// WriteCSV will write the connection in CSV format to the passed writer.
// An error will be returned if one occurs, otherwise error will be nil.
func (c Connection) WriteCSV(w *csv.Writer) error {
	return w.Write([]string{
		fmt.Sprintf("%d", c.Id),
		c.Nickname,
		c.User,
		c.Host,
		c.Description,
		c.Args,
		c.Identity,
		c.Command,
	})
}

// WriteJSON will write the connection in JSON format to the passed writer.
// An error will be returned if one occurs, otherwise error will be nil.
func (c Connection) WriteJSON(w io.Writer) error {
	j, err := json.Marshal(c)

	if err != nil {
		return err
	}

	_, err = w.Write(j)

	return err
}

// WriteLineLong writes a list format, single-line string to the passed
// writer interface. This is most likely used in listing connections.
// This func will write all connection properties.
// An error will be returned if one occurs, otherwise error will be nil.
func (c Connection) WriteLineLong(w io.Writer) error {
	_, err := fmt.Fprintf(w, "%-*d %s %s %s %s %s %s %s\n",
		ListViewColumnWidths["id"], c.Id,
		misc.StringTrimmer(c.Nickname, ListViewColumnWidths["nickname"]),
		misc.StringTrimmer(c.User, ListViewColumnWidths["user"]),
		misc.StringTrimmer(c.Host, ListViewColumnWidths["host"]),
		misc.StringTrimmer(c.Description, ListViewColumnWidths["description"]),
		misc.StringTrimmer(c.Args, ListViewColumnWidths["args"]),
		misc.StringTrimmer(c.Identity, ListViewColumnWidths["identity"]),
		misc.StringTrimmer(c.Command, ListViewColumnWidths["command"]),
	)

	return err
}

// WriteLineShort writes a list format, single-line string to the passed
// writer interface. This is most likely used in listing connections.
// This func will write only some connection properties.
// An error will be returned if one occurs, otherwise error will be nil.
func (c Connection) WriteLineShort(w io.Writer) error {
	_, err := fmt.Fprintf(w, "%-*d %s %s %s %s\n",
		ListViewColumnWidths["id"], c.Id,
		misc.StringTrimmer(c.Nickname, ListViewColumnWidths["nickname"]),
		misc.StringTrimmer(c.User, ListViewColumnWidths["user"]),
		misc.StringTrimmer(c.Host, ListViewColumnWidths["host"]),
		misc.StringTrimmer(c.Description, ListViewColumnWidths["description"]),
	)

	return err
}

// String returns a string containing the connection nickname and ID, in the
// format: "nickname (id)".
func (c Connection) String() string {
	return fmt.Sprintf("%s (%d)", c.Nickname, c.Id)
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

	// Validate connection properties
	err := c.Validate()

	// In this case, we want a non-zero connection ID, so err must be nil before
	// continuing
	if err != nil {
		return err
	}

	// Does the ID exist?
	exists, err := c.db.Exists(c.Id)

	if err != nil {
		return err
	}

	if !exists {
		return ErrIdNotExist
	}

	// Try updating the connection
	_, err = c.db.connection.Exec(`
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
		sqlNullableInt64(c.Id),
		sqlNullableString(c.Nickname),
		sqlNullableString(c.Host),
		sqlNullableString(c.User),
		sqlNullableString(c.Description),
		sqlNullableString(c.Args),
		sqlNullableString(c.Identity),
		sqlNullableString(c.Command),
	)

	return err
}

// Validate runs checks against Connection properties. This should be run
// before performing write operations against the database, as its purpose is
// to catch potentially fix-able errors before making SQL angry.

// Currently, the only check performed is whether the nickname is in a valid
// format (ex. starts with a letter). Additional checks may be added in the
// future.
func (c Connection) Validate() error {
	// Validate Nickname
	if c.Nickname == "" {
		return ErrConnNoNickname
	}

	if err := ValidateNickname(c.Nickname); err != nil {
		return err
	}

	// Validate Host
	if c.Host == "" {
		return ErrConnNoHost
	}
	// Validate User
	// Validate Description
	// Validate Args
	// Validate Identity
	// Validate Command

	// Validate Id
	// This needs to be the last test, as non-zero connection IDs are not catastrophic
	if c.Id < 0 {
		return ErrInvalidId
	}

	if c.Id == 0 {
		return ErrConnIdZero
	}

	return nil
}
