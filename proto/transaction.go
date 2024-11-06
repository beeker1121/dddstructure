package proto

// Transaction defines a transaction.
type Transaction struct {
	ID             uint
	UserID         uint
	Type           string
	CardType       string
	AmountCaptured uint
	InvoiceID      uint
	Status         string
}

// TransactionPaymentMethodCard defines the card payment method.
type TransactionPaymentMethodCard struct {
	Number         string
	ExpirationDate string
}

// TransactionPaymentMethod defines the transaction payment method.
type TransactionPaymentMethod struct {
	Card *TransactionPaymentMethodCard
}

// TransactionProcessParams defines the transaction process parameters.
type TransactionProcessParams struct {
	ID            uint
	UserID        uint
	Type          string
	Amount        uint
	PaymentMethod TransactionPaymentMethod
	InvoiceID     uint
}
