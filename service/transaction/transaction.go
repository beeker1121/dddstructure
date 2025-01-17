package transaction

import (
	"log/slog"
	"strings"

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
	logger   *slog.Logger
}

// SetServices sets the services interface.
func (s *Service) SetServices(services *interfaces.Service) {
	s.services = services
}

// New creates a new service.
func New(s *storage.Storage, l *slog.Logger) *Service {
	return &Service{
		storage: s,
		logger:  l,
	}
}

// Process handles processing a transaction.
func (s *Service) Process(params *proto.TransactionProcessParams) (*proto.Transaction, error) {
	// Validate parameters.
	if err := s.ValidateProcessParams(params); err != nil {
		return nil, err
	}

	// Handle ID.
	if params.ID == 0 {
		params.ID = idCounter
		idCounter++
	}

	// Get card type.
	cardType := "unknown"
	if params.PaymentMethod.Card != nil {
		if strings.HasPrefix(params.PaymentMethod.Card.Number, "411111") {
			cardType = "visa"
		}
	}

	// Create a transaction.
	storaget, err := s.storage.Transaction.Create(&transaction.Transaction{
		ID:             params.ID,
		UserID:         params.UserID,
		Type:           params.Type,
		CardType:       cardType,
		AmountCaptured: params.Amount,
		InvoiceID:      params.InvoiceID,
		Status:         "approved",
	})
	if err != nil {
		s.logger.Error("storage.Transaction.Create() error",
			slog.Any("error", err))
		return nil, err
	}

	// Update an invoice.
	if params.Type == "refund" {
		// Get the invoice.
		servicei, err := s.services.Invoice.GetByID(params.InvoiceID)
		if err != nil {
			return nil, err
		}

		// Change amounts and status.
		servicei.AmountDue += storaget.AmountCaptured
		servicei.AmountPaid -= storaget.AmountCaptured
		servicei.Status = "pending"

		if _, err := s.services.Invoice.UpdateForTransaction(&proto.InvoiceUpdateForTransactionParams{
			ID:         &servicei.ID,
			AmountDue:  &servicei.AmountDue,
			AmountPaid: &servicei.AmountPaid,
			Status:     &servicei.Status,
		}); err != nil {
			s.logger.Error("storage.Invoice.UpdateForTransaction() error",
				slog.Any("error", err))
			return nil, err
		}
	}

	ret := &proto.Transaction{
		ID:             storaget.ID,
		UserID:         storaget.UserID,
		Type:           storaget.Type,
		CardType:       storaget.CardType,
		AmountCaptured: storaget.AmountCaptured,
		InvoiceID:      storaget.InvoiceID,
		Status:         storaget.Status,
	}

	return ret, nil
}
