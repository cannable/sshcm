package cdb

import (
	"database/sql"
	"fmt"
)

type ConnectionProperty struct {
	Name  string
	Value string
}

func (p ConnectionProperty) IsEmpty() bool {
	return len(p.Value) < 1
}

func (p ConnectionProperty) Validate() error {
	if IsValidProperty(p.Name) {
		return nil
	}
	return ErrInvalidConnectionProperty
}

func (p ConnectionProperty) SqlNullableValue() sql.NullString {
	return sql.NullString{String: p.Value, Valid: true}
}

func (p ConnectionProperty) String() string {
	return p.Value
}

func (p ConnectionProperty) StringTrimmed(len int) string {
	f := fmt.Sprintf("%-*s", len, p.Value)

	return f[:len]
}
