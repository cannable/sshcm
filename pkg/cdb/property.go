package cdb

import (
	"database/sql"
	"fmt"
)

type ConnectionProperty struct {
	parent *Connection
	Name   string
	Value  string
}

func (p *ConnectionProperty) EffectiveValue() (string, error) {
	// If a value is set, return that
	if !p.IsEmpty() {
		return p.Value, nil
	}

	// If a value isn't set, go grab the default value
	def, err := p.parent.db.GetDefault("binary")

	if err != nil {
		return "", err
	}

	// Even if the default is empty, return it
	return def, nil
}

func (p *ConnectionProperty) IsEmpty() bool {
	return len(p.Value) < 1
}

func (p *ConnectionProperty) Validate() error {
	if IsValidProperty(p.Name) {
		return nil
	}
	return ErrInvalidConnectionProperty
}

func (p *ConnectionProperty) SqlNullableValue() sql.NullString {
	return sql.NullString{String: p.Value, Valid: true}
}

func (p *ConnectionProperty) String() string {
	return p.Value
}

func (p *ConnectionProperty) StringTrimmed(len int) string {
	f := fmt.Sprintf("%-*s", len, p.Value)

	return f[:len]
}
