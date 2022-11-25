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

// HandleMerchantBilling handles billing all merchants.
func (s *Service) HandleMerchantBilling() error {
	// Get billing information for all merchants.
	mad, err := s.c.Billing.GetMerchantAmountsDue()
	if err != nil {
		return err
	}
	for _, v := range mad {
		if err := s.c.Billing.AddAmountPaid(v.ID, 100); err != nil {
			return err
		}
	}

	return nil
}

// MerchantAmountsDue defines defines the merchant amounts due.
type MerchantAmountDue struct {
	ID           uint
	MerchantID   uint
	MerchantName string
	AmountDue    uint
}

// GetMerchantAmountsDue gets the merchant amounts due.
func (s *Service) GetMerchantAmountsDue() ([]*MerchantAmountDue, error) {
	// Get user by ID.
	mad, err := s.c.Billing.GetMerchantAmountsDue()
	if err != nil {
		return nil, err
	}

	var servicemad []*MerchantAmountDue
	for _, v := range mad {
		mapped := &MerchantAmountDue{
			ID:           v.ID,
			MerchantID:   v.MerchantID,
			MerchantName: v.MerchantName,
			AmountDue:    v.AmountDue,
		}
		servicemad = append(servicemad, mapped)
	}

	return servicemad, nil
}

// AddAmountPaid deducts the given amount from the amount due for the given
// merchant.
func (s *Service) AddAmountPaid(accountingID, amount uint) error {
	if err := s.c.Billing.AddAmountPaid(accountingID, amount); err != nil {
		return err
	}

	return nil
}

// AddAmountDue creates a new core accounting entry for the given merchant and
// user with the given amount.
func (s *Service) AddAmountDue(merchantID, userID, amount uint) error {
	if err := s.c.Billing.AddAmountDue(merchantID, userID, amount); err != nil {
		return err
	}

	return nil
}
