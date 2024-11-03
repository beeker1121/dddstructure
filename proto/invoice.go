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

// InvoiceGetParams defines the invoice get parameters.
type InvoiceGetParams struct {
	ID     *uint
	UserID *uint
	Status *string
	Offset uint
	Limit  uint
}

// InvoiceUpdateParams defines the invoice update parameters.
type InvoiceUpdateParams struct {
	ID         *uint
	UserID     *uint
	BillTo     *string
	PayTo      *string
	AmountDue  *uint
	AmountPaid *uint
	Status     *string
}

// InvoicePayParams defines the invoice pay parameters.
type InvoicePayParams struct {
	Amount uint
}
