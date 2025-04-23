package cdb

import (
	"database/sql"
)

func (conndb *ConnectionDB) GetDefault(name string) (string, error) {
	var def sql.NullString

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

func (conndb *ConnectionDB) SetDefault(setting string, value string) error {
	if !IsValidDefault(setting) {
		return ErrInvalidDefault
	}

	// Try updating the connection
	_, err := conndb.connection.Exec(`
		UPDATE defaults SET
			value = $2
		WHERE setting = $1
		`,
		setting,
		sqlNullableString(value),
	)

	return err
}

func (conndb *ConnectionDB) GetEffectiveValue(value string, setting string) (string, error) {
	// If a value is set, return that
	if len(value) > 0 {
		return value, nil
	}

	// If a value isn't set, go grab the default value
	def, err := conndb.GetDefault(setting)

	if err != nil {
		return "", err
	}

	// Even if the default is empty, return it
	return def, nil
}
