package merchant

// Database defines the merchant database interface.
type Database interface {
	Create(params *CreateParams) (*Merchant, error)
	GetByID(id uint) (*Merchant, error)
}

// Merchant defines the merchant.
type Merchant struct {
	ID    uint
	Name  string
	Email string
}

// CreateParams defines the create parameters.
type CreateParams struct {
	ID    uint
	Name  string
	Email string
}
