package errors

import "errors"

var (
	// ErrUserNotFound is returned when a user could not be found.
	ErrUserNotFound = errors.New("user not found")
)
