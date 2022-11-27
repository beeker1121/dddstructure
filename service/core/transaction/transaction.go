package transaction

import (
	"dddstructure/storage"
	"dddstructure/storage/transaction"
)

// idIncrementer handles incrementing the transaction IDs.
var idIncrementer uint = 1

// Core defines the core transaction service.
type Core struct {
	s *storage.Storage
}

// New creates a new core.
func New(s *storage.Storage) *Core {
	return &Core{
		s: s,
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

// CreateParams defines the Create parameters.
type CreateParams struct {
	ID             uint
	MerchantID     uint
	Type           string
	CardType       string
	AmountCaptured uint
}

// Create creates a new transaction.
func (c *Core) Create(params *CreateParams) (*Transaction, error) {
	// Handle creating an ID if one is not present.
	if params.ID == 0 {
		params.ID = idIncrementer
		idIncrementer++
	}

	// Create a transaction.
	t, err := c.s.Transaction.Create(&transaction.CreateParams{
		ID:             params.ID,
		MerchantID:     params.MerchantID,
		Type:           params.Type,
		CardType:       params.CardType,
		AmountCaptured: params.AmountCaptured,
	})
	if err != nil {
		return nil, err
	}

	// Map to core type.
	corem := &Transaction{
		ID:             t.ID,
		MerchantID:     t.MerchantID,
		Type:           t.Type,
		CardType:       t.CardType,
		AmountCaptured: t.AmountCaptured,
	}

	return corem, nil
}

// GetByID gets a transaction by the given ID.
func (c *Core) GetByID(id uint) (*Transaction, error) {
	// Get transaction by ID.
	t, err := c.s.Transaction.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Map to core type.
	coret := &Transaction{
		ID:         t.ID,
		MerchantID: t.MerchantID,
	}

	return coret, nil
}
