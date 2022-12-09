package invoice

import (
	"dddstructure/dep"
	"dddstructure/proto"
	"dddstructure/storage"
	"dddstructure/storage/invoice"
)

// idCounter handles increasing the ID.
var idCounter uint = 0

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

	// Create an invoice.
	inv, err := s.s.Invoice.Create(&invoice.Invoice{
		ID:         i.ID,
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
	i, err := s.s.Invoice.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Map to service type.
	servicei := &proto.Invoice{
		ID:         i.ID,
		BillTo:     i.BillTo,
		PayTo:      i.PayTo,
		AmountDue:  i.AmountDue,
		AmountPaid: i.AmountPaid,
		Status:     i.Status,
	}

	return servicei, nil
}

func (s *Service) Update(i *proto.Invoice) error {
	err := s.s.Invoice.Update(&invoice.Invoice{
		ID:         i.ID,
		BillTo:     i.BillTo,
		PayTo:      i.PayTo,
		AmountDue:  i.AmountDue,
		AmountPaid: i.AmountPaid,
		Status:     i.Status,
	})
	if err != nil {
		return err
	}

	// Pay the invoice using dependencies package.
	_, err = dep.Transaction.Process(&proto.Transaction{
		MerchantID:     1,
		Type:           "refund",
		CardType:       "visa",
		AmountCaptured: 100,
		InvoiceID:      0,
	})
	if err != nil {
		return err
	}

	return nil
}

// Pay handles paying an invoice.
func (s *Service) Pay(id uint) (*proto.Invoice, error) {
	// Need to call top level transaction.Process() service... yet
	// transaction.Process() is a 'top level' service, just like this
	// Pay() method is, and a top level service should not import
	// another top level service - the top level service should only
	// import core services.
	//
	// If we were to move transaction.Process() to a 'core level' service,
	// then that brings in another issue, since transaction.Process() will
	// need to call the core level processors service to actually process
	// the transaction - but a core level service should not import other
	// core level services.

	// Get the invoice.
	inv, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Pay the invoice using dependencies package.
	_, err = dep.Transaction.Process(&proto.Transaction{
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
	inv.AmountPaid = inv.AmountDue
	inv.AmountDue -= inv.AmountDue
	inv.Status = "paid"
	err = s.Update(inv)
	if err != nil {
		return nil, err
	}

	return inv, nil
}
