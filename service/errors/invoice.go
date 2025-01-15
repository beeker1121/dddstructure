package errors

import "errors"

var (
	// ErrInvoiceNotFound is returned when an invoice could not be found.
	ErrInvoiceNotFound = errors.New("invoice not found")

	// ErrInvoicePaymentMethodRequired is returned when no payment methods are
	// passed in.
	ErrInvoicePaymentMethodRequired = errors.New("at least one payment method is required")

	// ErrInvoicePaymentMethodInvalid is returned when at least one payment
	// method is invalid.
	ErrInvoicePaymentMethodInvalid = errors.New("invalid payment method, must be either 'card' or 'ach'")

	// ErrInvoiceAmountDueLimit is returned when the invoice amount due is over
	// the max limit.
	ErrInvoiceAmountDueLimit = errors.New("amount due is over limit")

	// ErrInvoiceCalculatingAmounts is returned when there was an error
	// calculating the invoice amounts.
	ErrInvoiceCalculatingAmounts = errors.New("error calculating invoice amounts")

	// ErrInvoiceStatusNotPending is returned when an invoice is trying to be
	// paid and is not in pending status.
	ErrInvoiceStatusNotPending = errors.New("invoice is not in pending status")
)
