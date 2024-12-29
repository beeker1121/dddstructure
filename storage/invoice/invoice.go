package invoice

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
	FirstName string
	LastName  string
}

// PayTo defines the payee information.
type PayTo struct {
	FirstName string
	LastName  string
}

// Invoice defines an invoice.
type Invoice struct {
	ID         uint
	UserID     uint
	BillTo     BillTo
	PayTo      PayTo
	AmountDue  uint
	AmountPaid uint
	Status     string
}

// GetParams defines the get parameters.
type GetParams struct {
	ID     *uint
	UserID *uint
	Status *string
	Offset uint
	Limit  uint
}
