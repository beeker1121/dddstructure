package errors

import "errors"

var (
	// ErrInvoiceNotFound is returned when an invoice could not be found.
	ErrInvoiceNotFound = errors.New("invoice not found")

	// ErrInvoiceAmountDueLimit is returned when the invoice amount due is over
	// the max limit.
	ErrInvoiceAmountDueLimit = errors.New("amount due is over limit")

	// ErrInvoiceCalculatingAmounts is returned when there was an error
	// calculating the invoice amounts.
	ErrInvoiceCalculatingAmounts = errors.New("error calculating invoice amounts")
)
