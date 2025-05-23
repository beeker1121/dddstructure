package auth

import "errors"

var (
	// ErrUnauthorized is returned when there is an error during JWT
	// authorization.
	ErrUnauthorized = errors.New("could not find JWT or API key in Authorization header")

	// ErrJWTUnauthorized is returned when there is an error during JWT
	// authorization.
	ErrJWTUnauthorized = errors.New("unauthorized")
)
