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

// InvoiceCreateParams defines the invoice create parameters.
type InvoiceCreateParams struct {
	ID        uint
	UserID    uint
	BillTo    string
	PayTo     string
	AmountDue uint
}

// InvoiceUpdateParams defines the invoice update parameters.
type InvoiceUpdateParams struct {
	ID         uint
	UserID     uint
	BillTo     string
	PayTo      string
	AmountDue  *uint
	AmountPaid *uint
	Status     *string
}

// InvoicePayParams defines the invoice pay parameters.
type InvoicePayParams struct {
	Amount uint
}
