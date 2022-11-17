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
	// ErrInvalidJson = errors.New("invalid json").
	// ErrBadRequest  = errors.New("bad request").
	// ErrNotFound    = errors.New("not found").
)
