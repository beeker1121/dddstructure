package proto

// Transaction defines a transaction.
type Transaction struct {
	ID             uint
	UserID         uint
	Type           string
	CardType       string
	AmountCaptured uint
	InvoiceID      uint
}

// TransactionProcessParams defines the transaction process parameters.
type TransactionProcessParams struct {
	ID        uint
	UserID    uint
	Type      string
	CardType  string
	Amount    uint
	InvoiceID uint
}
