package cdb

import "errors"

var ErrAddNicknameExists = errors.New("Add: nickname exists")
var ErrAddNoNickname = errors.New("Add: no nickname")
var ErrCreateSchemaVerInvalid = errors.New("createDB: invalid schema version")
var ErrMarshallNoRows = errors.New("marshallConnection: sql query returned 0 rows")
var ErrOpenDbVersionNotRecognized = errors.New("Open: DB schema version not recognized")
var ErrPropertyInvalid = errors.New("property is invalid")
var ErrUpdateIdNotExist = errors.New("Update: id does not exist")
var ErrUpdateNoId = errors.New("Update: no id")
var ErrUpdateNoNickname = errors.New("Update: no nickname")
