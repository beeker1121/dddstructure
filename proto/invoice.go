package proto

type Invoice struct {
	ID         uint
	BillTo     string
	PayTo      string
	AmountDue  uint
	AmountPaid uint
}
