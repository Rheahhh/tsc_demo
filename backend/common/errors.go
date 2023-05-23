package common

import "errors"

// Error definitions
var (
	ErrCurrentUser      = errors.New("cannot get current user")
	ErrDatabaseNotExist = errors.New("history database does not exist")
	ErrDatabaseOpen     = errors.New("cannot open database")
	ErrDatabaseQuery    = errors.New("cannot query database")
	ErrRowScan          = errors.New("cannot scan row")
	ErrRowIteration     = errors.New("rows error")
)
