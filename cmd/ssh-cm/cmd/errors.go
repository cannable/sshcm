package cmd

import "errors"

var ErrInvalidDefault = errors.New("invalid default")
var ErrNicknameExists = errors.New("nickname already exists")
var ErrNoIdOrNickname = errors.New("no id or nickname specified")
