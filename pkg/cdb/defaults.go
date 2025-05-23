package cdb

import (
	"database/sql"
)

// GetDefault retrieves a program default property from the connection database.
//
// If the passed property name is not valid, ErrInvalidDefault will be returned.
func (conndb *ConnectionDB) GetDefault(name string) (string, error) {
	var def sql.NullString

	if !IsValidDefault(name) {
		return "", ErrInvalidDefault
	}

	// Get connection details from DB
	err := conndb.connection.QueryRow(`
		SELECT value
		FROM defaults
		WHERE setting = $1
	`, name).Scan(&def)

	if err != nil {
		return "", err
	}

	return def.String, nil
}

// SetDefault updates a program default property in the connection database.
//
// If the passed property name is not valid, ErrInvalidDefault will be returned.
func (conndb *ConnectionDB) SetDefault(name string, value string) error {
	if !IsValidDefault(name) {
		return ErrInvalidDefault
	}

	// Try updating the connection
	_, err := conndb.connection.Exec(`
		UPDATE defaults SET
			value = $2
		WHERE setting = $1
		`,
		name,
		sqlNullableString(value),
	)

	return err
}
