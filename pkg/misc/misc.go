package misc

import (
	"database/sql"
	"fmt"
)

func SqlNullableString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: true}
}

func SqlNullableInt64(i int64) sql.NullInt64 {
	return sql.NullInt64{Int64: i, Valid: true}
}

func StringTrimmer(s string, len int) string {
	f := fmt.Sprintf("%-*s", len, s)

	return f[:len]
}
