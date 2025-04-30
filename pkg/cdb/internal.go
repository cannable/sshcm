package cdb

import (
	"database/sql"
	"errors"
)

// rowToConnection converts a sql row to a Connection struct.
func rowToConnection(row *sql.Row) (Connection, error) {
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

// sqlNullableInt64 returns a sql.NullInt64 containing the passed int64.
func sqlNullableInt64(i int64) sql.NullInt64 {
	return sql.NullInt64{Int64: i, Valid: true}
}

// sqlNullableString returns a sql.NullString containing the passed String.
// This is a utility function to allow for nullifying database TEXT fields.
func sqlNullableString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: true}
}
