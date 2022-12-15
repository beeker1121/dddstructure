package transaction

import (
	"dddstructure/dep"
	"dddstructure/proto"
	"dddstructure/storage"
	"dddstructure/storage/transaction"
)

// idCounter handles increasing the ID.
var idCounter uint = 1

// Service defines the transaction service.
type Service struct {
	s *storage.Storage
}

// New creates a new service.
func New(s *storage.Storage) *Service {
	return &Service{
		s: s,
	}
}

// Process handles processing a transaction.
func (s *Service) Process(t *proto.Transaction) (*proto.Transaction, error) {
	// Handle ID.
	if t.ID == 0 {
		t.ID = idCounter
		idCounter++
	}

	// Get the processor.
	p := dep.Processor.GetProcessor(t)

	// Process the transaction.
	if err := p.Process(t); err != nil {
		return nil, err
	}

	// Save new transaction.
	_, err := s.s.Transaction.Create(&transaction.Transaction{
		ID:             t.ID,
		MerchantID:     t.MerchantID,
		Type:           t.Type,
		CardType:       t.CardType,
		AmountCaptured: t.AmountCaptured,
		InvoiceID:      t.InvoiceID,
	})
	if err != nil {
		return nil, err
	}

	// Update an invoice.
	if t.Type == "refund" {
		// Get the invoice.
		i, err := dep.Invoice.GetByID(t.InvoiceID)
		if err != nil {
			return nil, err
		}

		// Change amounts and status.
		i.AmountDue += t.AmountCaptured
		i.AmountPaid -= t.AmountCaptured
		i.Status = "pending"

		// Update the invoice.
		if err := dep.Invoice.Update(i); err != nil {
			return nil, err
		}
	}

	return t, nil
}
