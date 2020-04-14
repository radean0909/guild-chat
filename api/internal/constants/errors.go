package constants

import (
	"errors"
)

var (
	// ErrNotFound - standard not found error - 404
	ErrNotFound = errors.New("not found")
	// ErrBadRequest - standard bad request error - typically due to bad data - 400
	ErrBadRequest = errors.New("bad request")
)
