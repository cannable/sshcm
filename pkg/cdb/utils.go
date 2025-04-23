package cdb

import (
	"strconv"
	"strings"
	"unicode"
)

var validDefaults = [8]string{
	"args",
	"binary",
	"command",
	"identity",
	"user",
}

var validProperties = [8]string{
	"nickname",
	"host",
	"user",
	"description",
	"args",
	"identity",
	"command",
	"binary",
}

func IsValidDefault(name string) bool {
	for _, v := range validDefaults {
		if strings.Compare(name, v) == 0 {
			return true
		}
	}
	return false
}

func IsValidProperty(property string) bool {
	for _, v := range validProperties {
		if strings.Compare(property, v) == 0 {
			return true
		}
	}
	return false
}

func ValidateNickname(nickname string) error {
	firstChar := []rune(nickname)[0]

	if !unicode.IsLetter(firstChar) {
		return ErrNicknameLetter
	}

	return nil
}

func ValidateId(id string) error {
	_, err := strconv.Atoi(id)

	if err != nil {
		return ErrInvalidId
	}

	return nil
}
