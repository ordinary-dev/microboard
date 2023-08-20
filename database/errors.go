package database

import (
	"errors"
)

var (
	ErrNoRowsWereAffected = errors.New("no rows were affected")
)
