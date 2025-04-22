package cmd

import "errors"

var ErrNicknameExists = errors.New("nickname already exists")
var ErrNoIdOrNickname = errors.New("no id or nickname specified")
var ErrNicknameLetter = errors.New("nickname does not begin with a letter")
