package cdb

import "errors"

var ErrConnFromDbInvalid = errors.New("connection from DB is invalid")
var ErrConnIdZero = errors.New("connection id is zero")
var ErrConnNoDb = errors.New("connection does not have a parent db attached")
var ErrConnNoHost = errors.New("connection does not have a host attached")
var ErrConnNoId = errors.New("connection does not have an id attached")
var ErrConnNoNickname = errors.New("connection does not have a nickname attached")
var ErrConnectionNotFound = errors.New("connection not found")
var ErrDuplicateNickname = errors.New("duplicate nickname")
var ErrIdNotExist = errors.New("connection id does not exist")
var ErrInvalidConnectionProperty = errors.New("invalid connection property")
var ErrInvalidDefault = errors.New("invalid default")
var ErrInvalidId = errors.New("invalid id")
var ErrInvalidIdOrNickname = errors.New("invalid id or nickname")
var ErrInvalidNickname = errors.New("invalid nickname")
var ErrNickNameNotExist = errors.New("connection nickname does not exist")
var ErrNicknameLetter = errors.New("nickname does not begin with a letter")
var ErrPropertyInvalid = errors.New("property is invalid")
var ErrUnsupportedSqlDriver = errors.New("sql driver not supported")

// DB schema errors
var ErrSchemaVerInvalid = errors.New("conndb: invalid schema version")
var ErrSchemaUpgradeNeeded = errors.New("conndb: schema upgrade needed")
var ErrSchemaNoUpgrade = errors.New("conndb: schema too old and can't be upgraded")
var ErrSchemaTooNew = errors.New("conndb: schema version too new")
