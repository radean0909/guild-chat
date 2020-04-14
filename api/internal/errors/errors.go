package errors

import (
	stdErrors "errors"
)

var (
	// ErrNotFound - standard not found error - 404
	ErrNotFound = stdErrors.New("not found")
	// ErrBadRequest - standard bad request error - typically due to bad data - 400
	ErrBadRequest = stdErrors.New("bad request")
)
