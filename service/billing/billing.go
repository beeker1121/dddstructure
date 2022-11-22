package billing

import (
	"dddstructure/service/core"
)

// Service defines the user service.
type Service struct {
	c *core.Core
}

// New creates a new service.
func New(c *core.Core) *Service {
	return &Service{
		c: c,
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
func (s *Service) GetMerchantAmountsDue() ([]*MerchantAmountsDue, error) {
	// Get user by ID.
	mad, err := s.c.Billing.GetMerchantAmountsDue()
	if err != nil {
		return nil, err
	}

	var servicemad []*MerchantAmountsDue
	for _, v := range mad {
		mapped := &MerchantAmountsDue{
			ID:           v.ID,
			MerchantID:   v.MerchantID,
			MerchantName: v.MerchantName,
			AmountDue:    v.AmountDue,
		}
		servicemad = append(servicemad, mapped)
	}

	return servicemad, nil
}
