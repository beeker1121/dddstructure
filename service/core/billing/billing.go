package billing

import (
	"dddstructure/storage"
	"dddstructure/storage/accounting"
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
type MerchantAmountDue struct {
	ID           uint
	MerchantID   uint
	MerchantName string
	AmountDue    uint
}

// GetMerchantAmountsDue gets the merchant amounts due.
func (c *Core) GetMerchantAmountsDue() ([]*MerchantAmountDue, error) {
	// Get account entry by ID.
	mad, err := c.s.Billing.GetMerchantAmountsDue()
	if err != nil {
		return nil, err
	}

	var coremad []*MerchantAmountDue
	for _, v := range mad {
		mapped := &MerchantAmountDue{
			ID:           v.ID,
			MerchantID:   v.MerchantID,
			MerchantName: v.MerchantName,
			AmountDue:    v.AmountDue,
		}
		coremad = append(coremad, mapped)
	}

	return coremad, nil
}

// AddAmountPaid deducts the given amount from the amount due for the given
// merchant.
func (c *Core) AddAmountPaid(accountingID, amount uint) error {
	// Get the accounting entry.
	a, err := c.s.Accounting.GetByID(accountingID)
	if err != nil {
		return err
	}

	// Modify the amount due.
	a.AmountDue -= amount

	// Update the accounting entry.
	c.s.Accounting.UpdateByID(accountingID, &accounting.UpdateParams{
		MerchantID: a.MerchantID,
		UserID:     a.UserID,
		AmountDue:  a.AmountDue,
	})

	return nil
}

// AddAmountDue creates a new core accounting entry for the given merchant and
// user with the given amount.
func (c *Core) AddAmountDue(merchantID, userID, amount uint) error {
	// Create a new accounting entry.
	_, err := c.s.Accounting.Create(&accounting.CreateParams{
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
