package invoice

import "time"

// Database defines the invoice database interface.
type Database interface {
	Create(i *Invoice) (*Invoice, error)
	Get(params *GetParams) ([]*Invoice, error)
	GetCount(params *GetParams) (uint, error)
	GetByID(id uint) (*Invoice, error)
	Update(i *Invoice) (*Invoice, error)
}

// BillTo defines the billing information.
type BillTo struct {
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

// PayTo defines the payee information.
type PayTo struct {
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

type LineItem struct {
	Name        string
	Description string
	Quantity    uint
	Price       uint
	Subtotal    uint
}

// Invoice defines an invoice.
type Invoice struct {
	ID            uint
	UserID        uint
	InvoiceNumber string
	PONumber      string
	Currency      string
	DueDate       time.Time
	Message       string
	BillTo        BillTo
	PayTo         PayTo
	LineItems     []LineItem
	TaxRate       string
	AmountDue     uint
	AmountPaid    uint
	Status        string
}

// GetParams defines the get parameters.
type GetParams struct {
	ID     *uint
	UserID *uint
	Status *string
	Offset uint
	Limit  uint
}
