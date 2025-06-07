package cdb

import (
	"database/sql"
)

// sqlNullableInt64 returns a sql.NullInt64 containing the passed int64.
func sqlNullableInt64(i int64) sql.NullInt64 {
	return sql.NullInt64{Int64: i, Valid: true}
}

// sqlNullableString returns a sql.NullString containing the passed String.
// This is a utility function to allow for nullifying database TEXT fields.
func sqlNullableString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: true}
}
