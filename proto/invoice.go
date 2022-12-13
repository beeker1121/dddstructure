package proto

type Invoice struct {
	ID         uint
	MerchantID uint
	BillTo     string
	PayTo      string
	AmountDue  uint
	AmountPaid uint
	Status     string
}
