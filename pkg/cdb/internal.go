package cdb

import (
	"database/sql"
	"errors"
)

func marshallConnection(row *sql.Row) (Connection, error) {
	c := NewConnection()

	var id sql.NullInt64
	var nickname, host, user, description, args, identity, command sql.NullString

	// Get connection details from DB
	err := row.Scan(&id, &nickname, &host, &user, &description, &args, &identity, &command)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c, ErrConnectionNotFound
		}
		return c, err
	}

	// If we're missing an id, nickname, or host name, something's very wrong
	if !(id.Valid || nickname.Valid || host.Valid) {
		return c, err
	}

	// Set connection ID, nickname, and host
	c.Id = id.Int64
	c.Nickname = nickname.String
	c.Host = host.String

	// The remaining bits from the DB are optional

	// If the DB has a user name for this connection, use it
	if user.Valid {
		c.User = user.String
	}

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

	return c, err
}

func sqlNullableString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}

	return sql.NullString{String: s, Valid: true}
}
