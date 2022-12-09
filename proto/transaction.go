package proto

type Transaction struct {
	ID             uint
	MerchantID     uint
	Type           string
	CardType       string
	AmountCaptured uint
}
