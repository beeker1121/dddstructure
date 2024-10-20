package proto

// Invoice defines an invoice.
type Invoice struct {
	ID         uint
	UserID     uint
	BillTo     string
	PayTo      string
	AmountDue  uint
	AmountPaid uint
	Status     string
}
