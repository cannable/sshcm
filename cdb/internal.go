package cdb

import (
	"database/sql"
	"errors"
	"strings"
)

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

func sqlNullableString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}

	return sql.NullString{String: s, Valid: true}
}
