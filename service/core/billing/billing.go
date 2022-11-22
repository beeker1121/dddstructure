package billing

import (
	"dddstructure/storage"
)

// Core defines the core accounting service.
type Core struct {
	s *storage.Storage
}

// New creates a new core.
func New(s *storage.Storage) *Core {
	return &Core{
		s: s,
	}
}

// MerchantAmountsDue defines defines the merchant amounts due.
type MerchantAmountsDue struct {
	ID           uint
	MerchantID   uint
	MerchantName string
	AmountDue    uint
}

// GetMerchantAmountsDue gets the merchant amounts due.
func (c *Core) GetMerchantAmountsDue() ([]*MerchantAmountsDue, error) {
	// Get account entry by ID.
	mad, err := c.s.Billing.GetMerchantAmountsDue()
	if err != nil {
		return nil, err
	}

	var coremad []*MerchantAmountsDue
	for _, v := range mad {
		mapped := &MerchantAmountsDue{
			ID:           v.ID,
			MerchantID:   v.MerchantID,
			MerchantName: v.MerchantName,
			AmountDue:    v.AmountDue,
		}
		coremad = append(coremad, mapped)
	}

	return coremad, nil
}
