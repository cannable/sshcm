package cdb

import "errors"

var ErrConnNoDb = errors.New("connection does not have a parent db attached")
var ErrConnNoId = errors.New("no id in connection struct")
var ErrConnNoNickname = errors.New("no nickname in connection struct")
var ErrDbVersionNotRecognized = errors.New("connection db schema version not recognized")
var ErrDuplicateNickname = errors.New("duplicate nickname")
var ErrIdNotExist = errors.New("connection id does not exist in db")
var ErrConnectionNotFound = errors.New("connection not found")
var ErrPropertyInvalid = errors.New("property is invalid")
var ErrSchemaVerInvalid = errors.New("invalid schema version")
