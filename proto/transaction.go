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
