package transaction

import (
	"dddstructure/service/core"
	"dddstructure/service/core/transaction"
)

// Service defines the transaction service.
type Service struct {
	c *core.Core
}

// New creates a new service.
func New(c *core.Core) *Service {
	return &Service{
		c: c,
	}
}

// Transaction defines a transaction.
type Transaction struct {
	ID             uint
	MerchantID     uint
	Type           string
	CardType       string
	AmountCaptured uint
}

// PaymentMethod defines a payment method.
type PaymentMethod struct {
	Card *Card
}

// Card defines a card.
type Card struct {
	Number         string
	ExpirationDate string
}

// ProcessParams defines the Process parameters.
type ProcessParams struct {
	ID            uint
	MerchantID    uint
	Type          string
	PaymentMethod PaymentMethod
	Amount        uint
}

// Process handles processing a transaction.
func (s *Service) Process(params *ProcessParams) (*Transaction, error) {
	// Get merchant.
	m, err := s.c.Merchant.GetByID(1)
	if err != nil {
		return nil, err
	}

	// Map to core processor service.

	// Process and get response.

	// Create params for new core transaction.
	createParams := &transaction.CreateParams{
		ID:         params.ID,
		MerchantID: params.MerchantID,
		Type:       params.Type,
	}

	// Fake processing transaction.
	if params.PaymentMethod.Card != nil {
		createParams.CardType = "visa"
	}
	var responseCode uint = 100
	switch responseCode {
	case 100:
		if params.Type == "capture" {
			createParams.AmountCaptured = params.Amount
		}
	}

	// Save new transaction.
	coret, err := s.c.Transaction.Create(createParams)
	if err != nil {
		return nil, err
	}

	// Handle updating invoices.
	if m.HasPermission("transactionModifiesInvoice") {
		// Get the invoice.
		i, err := s.c.Invoice.GetByID(1)
		if err != nil {
			return nil, err
		}

		switch coret.Type {
		case "capture":
			// Update invoice using core service.
			i.AmountDue -= params.Amount
			i.AmountPaid += params.Amount

			if err := s.c.Invoice.Update(i); err != nil {
				return nil, err
			}
		}
	}

	// Fake map response to transaction.
	t := &Transaction{
		ID:             coret.ID,
		MerchantID:     coret.MerchantID,
		Type:           coret.Type,
		CardType:       coret.CardType,
		AmountCaptured: coret.AmountCaptured,
	}

	return t, nil
}
