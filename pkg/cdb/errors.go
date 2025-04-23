package cdb

import "errors"

var ErrConnNoDb = errors.New("connection does not have a parent db attached")
var ErrConnNoId = errors.New("no id in connection struct")
var ErrConnNoNickname = errors.New("no nickname in connection struct")
var ErrConnectionNotFound = errors.New("connection not found")
var ErrDbVersionNotRecognized = errors.New("connection db schema version not recognized")
var ErrDuplicateNickname = errors.New("duplicate nickname")
var ErrIdNotExist = errors.New("connection id does not exist in db")
var ErrInvalidDefault = errors.New("invalid default")
var ErrInvalidId = errors.New("invalid id")
var ErrNicknameLetter = errors.New("nickname does not begin with a letter")
var ErrPropertyInvalid = errors.New("property is invalid")
var ErrSchemaVerInvalid = errors.New("invalid schema version")
