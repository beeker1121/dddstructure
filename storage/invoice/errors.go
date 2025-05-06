package invoice

import "errors"

var (
	// ErrInvoiceNotFound is returned when an invoice could not be found.
	ErrInvoiceNotFound = errors.New("invoice not found")
)
