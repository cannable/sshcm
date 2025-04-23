package cdb

import (
	"strconv"
	"strings"
	"unicode"
)

var ValidDefaults = [5]string{
	"args",
	"binary",
	"command",
	"identity",
	"user",
}

var ValidProperties = [8]string{
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
	for _, v := range ValidDefaults {
		if strings.Compare(name, v) == 0 {
			return true
		}
	}
	return false
}

func IsValidProperty(property string) bool {
	for _, v := range ValidProperties {
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
