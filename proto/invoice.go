package proto

import "time"

// InvoicePaymentMethod defines an invoice payment method.
type InvoicePaymentMethod string

const (
	InvoicePaymentMethodCard InvoicePaymentMethod = "card"
	InvoicePaymentMethodACH  InvoicePaymentMethod = "ach"
)

// InvoiceBillTo defines the invoice billing information.
type InvoiceBillTo struct {
	FirstName    string
	LastName     string
	Company      string
	AddressLine1 string
	AddressLine2 string
	City         string
	State        string
	PostalCode   string
	Country      string
	Email        string
	Phone        string
}

// InvoicePayTo defines the invoice payee information.
type InvoicePayTo struct {
	FirstName    string
	LastName     string
	Company      string
	AddressLine1 string
	AddressLine2 string
	City         string
	State        string
	PostalCode   string
	Country      string
	Email        string
	Phone        string
}

// InvoiceLineItem defines an invoice line item.
type InvoiceLineItem struct {
	Name        string
	Description string
	Quantity    uint
	Price       uint
}

// Invoice defines an invoice.
type Invoice struct {
	ID             uint
	UserID         uint
	InvoiceNumber  string
	PONumber       string
	Currency       string
	DueDate        time.Time
	Message        string
	BillTo         InvoiceBillTo
	PayTo          InvoicePayTo
	LineItems      []InvoiceLineItem
	PaymentMethods []InvoicePaymentMethod
	TaxRate        string
	AmountDue      uint
	AmountPaid     uint
	Status         string
	CreatedAt      time.Time
}

// InvoiceCreateParams defines the invoice create parameters.
type InvoiceCreateParams struct {
	ID             uint
	UserID         uint
	InvoiceNumber  string
	PONumber       string
	Currency       string
	DueDate        time.Time
	Message        string
	BillTo         InvoiceBillTo
	PayTo          InvoicePayTo
	LineItems      []InvoiceLineItem
	PaymentMethods []InvoicePaymentMethod
	TaxRate        string
}

// InvoiceGetParamsCreatedAt defines a created at datetime range.
type InvoiceGetParamsCreatedAt struct {
	StartDate *time.Time
	EndDate   *time.Time
}

// InvoiceGetParams defines the invoice get parameters.
type InvoiceGetParams struct {
	ID        *uint
	UserID    *uint
	Status    *string
	CreatedAt *InvoiceGetParamsCreatedAt
	Offset    uint
	Limit     uint
}

// InvoiceBillToUpdate defines the invoice billing information for update.
type InvoiceBillToUpdate struct {
	FirstName    *string
	LastName     *string
	Company      *string
	AddressLine1 *string
	AddressLine2 *string
	City         *string
	State        *string
	PostalCode   *string
	Country      *string
	Email        *string
	Phone        *string
}

// InvoicePayToUpdate defines the invoice payee information for update.
type InvoicePayToUpdate struct {
	FirstName    *string
	LastName     *string
	Company      *string
	AddressLine1 *string
	AddressLine2 *string
	City         *string
	State        *string
	PostalCode   *string
	Country      *string
	Email        *string
	Phone        *string
}

// InvoiceUpdateParams defines the invoice update parameters.
type InvoiceUpdateParams struct {
	ID             *uint
	UserID         *uint
	InvoiceNumber  *string
	PONumber       *string
	Currency       *string
	DueDate        *time.Time
	Message        *string
	BillTo         *InvoiceBillToUpdate
	PayTo          *InvoicePayToUpdate
	LineItems      *[]InvoiceLineItem
	PaymentMethods *[]InvoicePaymentMethod
	TaxRate        *string
}

// InvoiceUpdateForTransactionParams defines the invoice update for transaction
// parameters.
type InvoiceUpdateForTransactionParams struct {
	ID         *uint
	AmountDue  *uint
	AmountPaid *uint
	Status     *string
}

// InvoicePayParams defines the invoice pay parameters.
type InvoicePayParams struct {
	Amount uint
}
