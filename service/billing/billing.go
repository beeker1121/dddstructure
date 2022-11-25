package billing

import (
	"dddstructure/service/core"
	"dddstructure/service/core/accounting"
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
		// Add amount paid.
		_, err = s.c.Accounting.UpdateByID(v.ID, &accounting.UpdateParams{
			MerchantID: v.MerchantID,
			UserID:     1,
			AmountDue:  0,
		})
		if err != nil {
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
	// Get the accounting entry.
	a, err := s.c.Accounting.GetByID(accountingID)
	if err != nil {
		return err
	}

	// Modify the amount due.
	a.AmountDue -= amount

	// Update the accounting entry.
	s.c.Accounting.UpdateByID(accountingID, &accounting.UpdateParams{
		MerchantID: a.MerchantID,
		UserID:     a.UserID,
		AmountDue:  a.AmountDue,
	})

	return nil
}

// AddAmountDue creates a new core accounting entry for the given merchant and
// user with the given amount.
func (s *Service) AddAmountDue(merchantID, userID, amount uint) error {
	// Create a new accounting entry.
	_, err := s.c.Accounting.Create(&accounting.CreateParams{
		ID:         1,
		MerchantID: merchantID,
		UserID:     userID,
		AmountDue:  amount,
	})
	if err != nil {
		return err
	}

	return nil
}
