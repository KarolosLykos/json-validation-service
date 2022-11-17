package exceptions

import (
	"errors"
)

var (
	ErrRecover              = errors.New("recovering from error")
	ErrConnectingToDatabase = errors.New("could not connect to database")
	ErrGetDB                = errors.New("could not get database")
	ErrCloseDB              = errors.New("could not close database connection")
	ErrInitializeDatabase   = errors.New("could not initialize database")
	ErrInvalidJSON          = errors.New("invalid json")
	ErrInternalServerError  = errors.New("internal server error")
	ErrNotFound             = errors.New("not found")
	ErrValidation           = errors.New("error validating the given json data, against the json-schema")
)
