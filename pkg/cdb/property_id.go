package cdb

import (
	"database/sql"
	"fmt"
)

type IdProperty struct {
	parent *Connection
	Name   string
	Value  int64
}

func (p *IdProperty) EffectiveValue() (int64, error) {
	return p.Value, nil
}

func (p *IdProperty) IsEmpty() bool {
	return p.Value < 1
}

func (p *IdProperty) Validate() error {
	return ValidateId(p.String())
}

func (p *IdProperty) SqlNullableValue() sql.NullInt64 {
	return sql.NullInt64{Int64: p.Value, Valid: true}
}

func (p *IdProperty) String() string {
	return fmt.Sprintf("%d", p.Value)
}

func (p *IdProperty) StringTrimmed(len int) string {
	f := fmt.Sprintf("%-*d", len, p.Value)

	return f[:len]
}
