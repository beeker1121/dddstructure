package transaction

// Database defines the transaction database interface.
type Database interface {
	Create(i *Transaction) (*Transaction, error)
	GetByID(id uint) (*Transaction, error)
}

// Transaction defines the transaction.
type Transaction struct {
	ID             uint
	MerchantID     uint
	Type           string
	ProcessorType  string
	CardType       string
	AmountCaptured uint
	InvoiceID      uint
}
