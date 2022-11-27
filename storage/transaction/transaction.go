package transaction

// Database defines the transaction database interface.
type Database interface {
	Create(params *CreateParams) (*Transaction, error)
	GetByID(id uint) (*Transaction, error)
}

// Transaction defines the transaction.
type Transaction struct {
	ID             uint
	MerchantID     uint
	Type           string
	CardType       string
	AmountCaptured uint
}

// CreateParams defines the create parameters.
type CreateParams struct {
	ID             uint
	MerchantID     uint
	Type           string
	CardType       string
	AmountCaptured uint
}
