package cdb

import (
	"database/sql"
	"fmt"
)

type NicknameProperty struct {
	parent *Connection
	Name   string
	Value  string
}

func (p *NicknameProperty) EffectiveValue() (string, error) {
	return p.Value, nil
}

func (p *NicknameProperty) IsEmpty() bool {
	return len(p.Value) < 1
}

func (p *NicknameProperty) Validate() error {
	return ValidateNickname(p.Value)
}

func (p *NicknameProperty) SqlNullableValue() sql.NullString {
	return sql.NullString{String: p.Value, Valid: true}
}

func (p *NicknameProperty) String() string {
	return p.Value
}

func (p *NicknameProperty) StringTrimmed(len int) string {
	f := fmt.Sprintf("%-*s", len, p.Value)

	return f[:len]
}
