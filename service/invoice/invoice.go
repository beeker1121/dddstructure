package invoice

import (
	"dddstructure/proto"
	"dddstructure/service/interfaces"
	"dddstructure/storage"
	"dddstructure/storage/invoice"
)

// idCounter handles increasing the ID.
var idCounter uint = 1

// Service defines the invoice service.
type Service struct {
	storage  *storage.Storage
	services *interfaces.Service
}

func (s *Service) SetServices(services *interfaces.Service) {
	s.services = services
}

// New creates a new service.
func New(s *storage.Storage) *Service {
	return &Service{
		storage: s,
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
	inv, err := s.storage.Invoice.Create(&invoice.Invoice{
		ID:         i.ID,
		MerchantID: i.MerchantID,
		BillTo:     i.BillTo,
		PayTo:      i.PayTo,
		AmountDue:  i.AmountDue,
		AmountPaid: i.AmountPaid,
		Status:     i.Status,
	})
	if err != nil {
		return nil, err
	}

	// Map to service type.
	servicei := &proto.Invoice{
		ID:         inv.ID,
		MerchantID: inv.MerchantID,
		BillTo:     inv.BillTo,
		PayTo:      inv.PayTo,
		AmountDue:  inv.AmountDue,
		AmountPaid: inv.AmountPaid,
		Status:     inv.Status,
	}

	return servicei, nil
}

// GetByID gets an invoice by the given ID.
func (s *Service) GetByID(id uint) (*proto.Invoice, error) {
	// Get invoice by ID.
	i, err := s.storage.Invoice.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Map to service type.
	servicei := &proto.Invoice{
		ID:         i.ID,
		MerchantID: i.MerchantID,
		BillTo:     i.BillTo,
		PayTo:      i.PayTo,
		AmountDue:  i.AmountDue,
		AmountPaid: i.AmountPaid,
		Status:     i.Status,
	}

	return servicei, nil
}

func (s *Service) Update(i *proto.Invoice) error {
	return s.storage.Invoice.Update(&invoice.Invoice{
		ID:         i.ID,
		MerchantID: i.MerchantID,
		BillTo:     i.BillTo,
		PayTo:      i.PayTo,
		AmountDue:  i.AmountDue,
		AmountPaid: i.AmountPaid,
		Status:     i.Status,
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
	t, err := s.services.Transaction.Process(&proto.Transaction{
		MerchantID:     inv.MerchantID,
		Type:           "capture",
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
