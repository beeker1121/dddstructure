package errors

import "errors"

var (
	// ErrUserNotFound is returned when a user could not be found.
	ErrUserNotFound = errors.New("user not found")

	// ErrUserEmailEmpty is returned when the email param is empty.
	ErrUserEmailEmpty = errors.New("email parameter is empty")

	// ErrUserEmailExists is returned when the email already exists.
	ErrUserEmailExists = errors.New("email already exists")

	// ErrUserPassword is returned when the password is in an invalid format.
	ErrUserPassword = errors.New("password must be at least 8 characters")

	// ErrUserInvalidLogin is returned when the email and/or password used with
	// login is invalid.
	ErrUserInvalidLogin = errors.New("email and/or password is invalid")
)
