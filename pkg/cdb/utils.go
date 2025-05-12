package cdb

import (
	"strconv"
	"strings"
	"unicode"
)

var ValidDefaults = [4]string{
	"args",
	"command",
	"identity",
	"user",
}

var ValidProperties = [7]string{
	"nickname",
	"host",
	"user",
	"description",
	"args",
	"identity",
	"command",
}

// IsValidDefault checks the passed default property name against a list of
// valid default properties. This is similar to IsValidProperty, except there
// are fewer defaults.
//
// Returns true if the name is valid, false otherwise.
func IsValidDefault(name string) bool {
	for _, v := range ValidDefaults {
		if strings.Compare(name, v) == 0 {
			return true
		}
	}
	return false
}

// IsValidProperty checks the passed connection property name against a list of
// valid properties.
//
// Returns true if the name is valid, false otherwise.
func IsValidProperty(property string) bool {
	for _, v := range ValidProperties {
		if strings.Compare(property, v) == 0 {
			return true
		}
	}
	return false
}

// ValidateNickname runs checks against the passed nickname string.
//
// If the tests pass and the nickname is valid, nil is returned.
// If any test fails, a relevant error is returned.
func ValidateNickname(nickname string) error {
	if len(nickname) < 1 {
		return ErrInvalidNickname
	}

	firstChar := []rune(nickname)[0]

	if !unicode.IsLetter(firstChar) {
		return ErrNicknameLetter
	}

	return nil
}

// ValidateId runs checks against the passed id as a string.
//
// If the tests pass and the id is valid, nil is returned.
// If any test fails, a relevant error is returned.
//
// This is implemented as a string, as it is intended to be used in
// circumstances where a string might contain an id, or what is expected to be
// an id. The primary situation where this occurs is validating command line
// arguments - because a user may choose to make a connection by id or
// nickname, this func is the frontline validation of that user input.
func ValidateId(id string) error {
	// ids must be a valid integer
	i, err := strconv.Atoi(id)

	if err != nil {
		return ErrInvalidId
	}

	// id numbers must start at 1
	if i < 1 {
		return ErrInvalidId
	}

	return nil
}

// IsValidIdOrNickname returns true if the passed string is a valid id or
// nickname. This is a smoke test meant to be used to simplify conditionals
// around validating command line arguments.
func IsValidIdOrNickname(s string) bool {
	// Determine if the passed string is a nickname or id
	if err := ValidateId(s); err == nil {
		// Got a valid id
		return true
	} else if err := ValidateNickname(s); err == nil {
		// Got a valid nickname
		return true
	}
	return false
}
