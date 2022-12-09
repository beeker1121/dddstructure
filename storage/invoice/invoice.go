package invoice

// Database defines the invoice database interface.
type Database interface {
	Create(i *Invoice) (*Invoice, error)
	GetByID(id uint) (*Invoice, error)
	Update(i *Invoice) error
}

// Invoice defines an invoice.
type Invoice struct {
	ID         uint
	MerchantID uint
	BillTo     string
	PayTo      string
	AmountDue  uint
	AmountPaid uint
	Status     string
}
