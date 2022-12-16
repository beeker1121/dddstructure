package invoice

import (
	"dddstructure/dep"
	"dddstructure/proto"
	"dddstructure/storage"
	"dddstructure/storage/invoice"
)

// idCounter handles increasing the ID.
var idCounter uint = 1

// Service defines the invoice service.
type Service struct {
	s *storage.Storage
}

// New creates a new service.
func New(s *storage.Storage) *Service {
	return &Service{
		s: s,
	}
}

// Create creates a new invoice.
func (s *Service) Create(i *proto.Invoice) (*proto.Invoice, error) {
	// Handle ID.
	if i.ID == 0 {
		i.ID = idCounter
		idCounter++
	}

	// Handle status.
	if i.Status == "" {
		i.Status = "pending"
	}

	// Create an invoice.
	inv, err := s.s.Invoice.Create(&invoice.Invoice{
		ID:            i.ID,
		MerchantID:    i.MerchantID,
		ProcessorType: i.ProcessorType,
		BillTo:        i.BillTo,
		PayTo:         i.PayTo,
		AmountDue:     i.AmountDue,
		AmountPaid:    i.AmountPaid,
		Status:        i.Status,
	})
	if err != nil {
		return nil, err
	}

	// Map to service type.
	servicei := &proto.Invoice{
		ID:            inv.ID,
		MerchantID:    inv.MerchantID,
		ProcessorType: inv.ProcessorType,
		BillTo:        inv.BillTo,
		PayTo:         inv.PayTo,
		AmountDue:     inv.AmountDue,
		AmountPaid:    inv.AmountPaid,
		Status:        inv.Status,
	}

	return servicei, nil
}

// GetByID gets an invoice by the given ID.
func (s *Service) GetByID(id uint) (*proto.Invoice, error) {
	// Get invoice by ID.
	i, err := s.s.Invoice.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Map to service type.
	servicei := &proto.Invoice{
		ID:            i.ID,
		MerchantID:    i.MerchantID,
		ProcessorType: i.ProcessorType,
		BillTo:        i.BillTo,
		PayTo:         i.PayTo,
		AmountDue:     i.AmountDue,
		AmountPaid:    i.AmountPaid,
		Status:        i.Status,
	}

	return servicei, nil
}

func (s *Service) Update(i *proto.Invoice) error {
	return s.s.Invoice.Update(&invoice.Invoice{
		ID:            i.ID,
		MerchantID:    i.MerchantID,
		ProcessorType: i.ProcessorType,
		BillTo:        i.BillTo,
		PayTo:         i.PayTo,
		AmountDue:     i.AmountDue,
		AmountPaid:    i.AmountPaid,
		Status:        i.Status,
	})
}

// Pay handles paying an invoice.
func (s *Service) Pay(id uint) (*proto.Invoice, error) {
	// Get the invoice.
	inv, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Pay the invoice using dependencies package.
	t, err := dep.Transaction.Process(&proto.Transaction{
		MerchantID:     inv.MerchantID,
		Type:           "capture",
		ProcessorType:  inv.ProcessorType,
		CardType:       "visa",
		AmountCaptured: inv.AmountDue,
		InvoiceID:      id,
	})
	if err != nil {
		return nil, err
	}

	// Update the invoice.
	inv.AmountPaid = t.AmountCaptured
	inv.AmountDue -= t.AmountCaptured
	inv.Status = "paid"
	err = s.Update(inv)
	if err != nil {
		return nil, err
	}

	return inv, nil
}
