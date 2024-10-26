package errors

import "errors"

var (
	// ErrTransactionNotFound is returned when a transaction could not be
	// found.
	ErrTransactionNotFound = errors.New("transaction not found")
)
