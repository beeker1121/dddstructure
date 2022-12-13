package merchant

// Database defines the merchant database interface.
type Database interface {
	Create(m *Merchant) (*Merchant, error)
	GetByID(id uint) (*Merchant, error)
}

// Merchant defines an merchant.
type Merchant struct {
	ID    uint
	Name  string
	Email string
}
