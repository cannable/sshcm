package cdb

import "errors"

var ErrConnNoDb = errors.New("connection does not have a parent db attached")
var ErrConnNoId = errors.New("connection does not have an id attached")
var ErrConnNoNickname = errors.New("connection does not have a nickname attached")
var ErrConnectionNotFound = errors.New("connection not found")
var ErrDbVersionNotRecognized = errors.New("connection db schema version not recognized")
var ErrDuplicateNickname = errors.New("duplicate nickname")
var ErrIdNotExist = errors.New("connection id does not exist")
var ErrInvalidConnectionProperty = errors.New("invalid connection property")
var ErrInvalidDefault = errors.New("invalid default")
var ErrInvalidId = errors.New("invalid id")
var ErrNickNameNotExist = errors.New("connection nickname does not exist")
var ErrNicknameLetter = errors.New("nickname does not begin with a letter")
var ErrPropertyInvalid = errors.New("property is invalid")
var ErrSchemaVerInvalid = errors.New("invalid schema version")
