package transaction

import (
	"dddstructure/proto"
	"dddstructure/service/interfaces"
	"dddstructure/storage"
	"dddstructure/storage/transaction"
)

// idCounter handles increasing the ID.
var idCounter uint = 1

// Service defines the transaction service.
type Service struct {
	storage  *storage.Storage
	services *interfaces.Service
}

// SetServices sets the services interface.
func (s *Service) SetServices(services *interfaces.Service) {
	s.services = services
}

// New creates a new service.
func New(s *storage.Storage) *Service {
	return &Service{
		storage: s,
	}
}

// Process handles processing a transaction.
func (s *Service) Process(t *proto.TransactionProcessParams) (*proto.Transaction, error) {
	// Handle ID.
	if t.ID == 0 {
		t.ID = idCounter
		idCounter++
	}

	// Save new transaction.
	st, err := s.storage.Transaction.Create(&transaction.Transaction{
		ID:             t.ID,
		UserID:         t.UserID,
		Type:           t.Type,
		CardType:       t.CardType,
		AmountCaptured: t.Amount,
		InvoiceID:      t.InvoiceID,
	})
	if err != nil {
		return nil, err
	}

	// Update an invoice.
	if t.Type == "refund" {
		// Get the invoice.
		i, err := s.services.Invoice.GetByID(t.InvoiceID)
		if err != nil {
			return nil, err
		}

		// Change amounts and status.
		i.AmountDue += st.AmountCaptured
		i.AmountPaid -= st.AmountCaptured
		i.Status = "pending"

		s.services.Invoice.Update(i)
	}

	ret := &proto.Transaction{
		ID:             st.ID,
		UserID:         st.UserID,
		Type:           st.Type,
		CardType:       st.CardType,
		AmountCaptured: st.AmountCaptured,
		InvoiceID:      st.InvoiceID,
	}

	return ret, nil
}
