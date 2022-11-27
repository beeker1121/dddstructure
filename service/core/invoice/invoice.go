package invoice

import (
	"dddstructure/storage"
	"dddstructure/storage/invoice"
)

// idIncrementer handles incrementing the invoice IDs.
var idIncrementer uint = 1

// Core defines the core invoice service.
type Core struct {
	s *storage.Storage
}

// New creates a new core.
func New(s *storage.Storage) *Core {
	return &Core{
		s: s,
	}
}

// Invoice defines an invoice.
type Invoice struct {
	ID         uint
	BillTo     string
	PayTo      string
	AmountDue  uint
	AmountPaid uint
}

// CreateParams defines the Create parameters.
type CreateParams struct {
	ID         uint
	BillTo     string
	PayTo      string
	AmountDue  uint
	AmountPaid uint
}

// Create creates a new invoice.
func (c *Core) Create(params *CreateParams) (*Invoice, error) {
	// Handle creating an ID if one is not present.
	if params.ID == 0 {
		params.ID = idIncrementer
		idIncrementer++
	}

	// Create an invoice.
	i, err := c.s.Invoice.Create(&invoice.CreateParams{
		ID:         params.ID,
		BillTo:     params.BillTo,
		PayTo:      params.PayTo,
		AmountDue:  params.AmountDue,
		AmountPaid: params.AmountPaid,
	})
	if err != nil {
		return nil, err
	}

	// Map to core type.
	corei := &Invoice{
		ID:         i.ID,
		BillTo:     i.BillTo,
		PayTo:      i.PayTo,
		AmountDue:  i.AmountDue,
		AmountPaid: i.AmountPaid,
	}

	return corei, nil
}

// GetByID gets an invoice by the given ID.
func (c *Core) GetByID(id uint) (*Invoice, error) {
	// Get invoice by ID.
	i, err := c.s.Invoice.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Map to core type.
	corei := &Invoice{
		ID:         i.ID,
		BillTo:     i.BillTo,
		PayTo:      i.PayTo,
		AmountDue:  i.AmountDue,
		AmountPaid: i.AmountPaid,
	}

	return corei, nil
}

// Update updates an invoice.
func (c *Core) Update(i *Invoice) error {
	// Get the invoice from storage.
	storeagei, err := c.s.Invoice.GetByID(i.ID)
	if err != nil {
		return err
	}

	// Update fields.
	storeagei.ID = i.ID
	storeagei.BillTo = i.BillTo
	storeagei.PayTo = i.PayTo
	storeagei.AmountDue = i.AmountDue
	storeagei.AmountPaid = i.AmountPaid

	// Update the invoice via storage.
	if err := c.s.Invoice.Update(storeagei); err != nil {
		return err
	}

	return nil
}
