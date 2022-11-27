package invoice

import (
	"dddstructure/service/core"
	"dddstructure/service/core/invoice"
)

// Service defines the invoice service.
type Service struct {
	c *core.Core
}

// New creates a new service.
func New(c *core.Core) *Service {
	return &Service{
		c: c,
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
func (s *Service) Create(params *CreateParams) (*Invoice, error) {
	// Create an invoice.
	i, err := s.c.Invoice.Create(&invoice.CreateParams{
		ID:         params.ID,
		BillTo:     params.BillTo,
		PayTo:      params.PayTo,
		AmountDue:  params.AmountDue,
		AmountPaid: params.AmountPaid,
	})
	if err != nil {
		return nil, err
	}

	// Map to service type.
	servicei := &Invoice{
		ID:         i.ID,
		BillTo:     i.BillTo,
		PayTo:      i.PayTo,
		AmountDue:  i.AmountDue,
		AmountPaid: i.AmountPaid,
	}

	return servicei, nil
}

// GetByID gets an invoice by the given ID.
func (s *Service) GetByID(id uint) (*Invoice, error) {
	// Get invoice by ID.
	i, err := s.c.Invoice.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Map to service type.
	servicei := &Invoice{
		ID:         i.ID,
		BillTo:     i.BillTo,
		PayTo:      i.PayTo,
		AmountDue:  i.AmountDue,
		AmountPaid: i.AmountPaid,
	}

	return servicei, nil
}
