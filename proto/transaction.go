package proto

type Transaction struct {
	ID             uint
	MerchantID     uint
	Type           string
	ProcessorType  string
	CardType       string
	AmountCaptured uint
	InvoiceID      uint
}
