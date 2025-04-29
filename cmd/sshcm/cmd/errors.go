package cmd

import "errors"

var ErrImportCSVInvalidColumn = errors.New("spurious column in import file")
var ErrImportFileNotFound = errors.New("import file does not exist")
var ErrInvalidDefault = errors.New("invalid default")
var ErrNicknameExists = errors.New("nickname already exists")
var ErrNoIdOrNickname = errors.New("no id or nickname specified")
