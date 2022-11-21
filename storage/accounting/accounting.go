package accounting

// Database defines the accounting database interface.
type Database interface {
	Create(params *CreateParams) (*Accounting, error)
	GetByID(id uint) (*Accounting, error)
}

// Accounting defines an accounting entry.
type Accounting struct {
	ID         uint
	MerchantID uint
	UserID     uint
	AmountDue  uint
}

// CreateParams defines the create parameters.
type CreateParams struct {
	ID         uint
	MerchantID uint
	UserID     uint
	AmountDue  uint
}
