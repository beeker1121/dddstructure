package invoice

// Database defines the invoice database interface.
type Database interface {
	Create(params *CreateParams) (*Invoice, error)
	GetByID(id uint) (*Invoice, error)
	Update(i *Invoice) error
}

// Invoice defines an invoice.
type Invoice struct {
	ID         uint
	BillTo     string
	PayTo      string
	AmountDue  uint
	AmountPaid uint
}

// CreateParams defines the create parameters.
type CreateParams struct {
	ID         uint
	BillTo     string
	PayTo      string
	AmountDue  uint
	AmountPaid uint
}
