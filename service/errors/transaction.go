package errors

import "errors"

var (
	// ErrTransactionNotFound is returned when a transaction could not be
	// found.
	ErrTransactionNotFound = errors.New("transaction not found")

	// ErrTransactionAmountLimit is returned when the transaction amount is
	// over the max limit.
	ErrTransactionAmountLimit = errors.New("transaction amount is over limit")
)
