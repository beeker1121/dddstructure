package accounting

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

// Accounting defines an accounting entry.
type Accounting struct {
	ID         uint
	MerchantID uint
	UserID     uint
	AmountDue  uint
}

// CreateParams defines the Create parameters.
type CreateParams struct {
	ID         uint
	MerchantID uint
	UserID     uint
	AmountDue  uint
}

// Create creates a new accounting entry.
func (c *Core) Create(params *CreateParams) (*Accounting, error) {
	// Create a new accounting entry.
	a, err := c.s.Accounting.Create(&accounting.CreateParams{
		ID:         params.ID,
		MerchantID: params.MerchantID,
		UserID:     params.UserID,
		AmountDue:  params.AmountDue,
	})
	if err != nil {
		return nil, err
	}

	// Map to core type.
	corea := &Accounting{
		ID:         a.ID,
		MerchantID: a.MerchantID,
		UserID:     a.UserID,
		AmountDue:  a.AmountDue,
	}

	return corea, nil
}

// GetByID gets an accounting by the given ID.
func (c *Core) GetByID(id uint) (*Accounting, error) {
	// Get account entry by ID.
	a, err := c.s.Accounting.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Map to core type.
	corea := &Accounting{
		ID:         a.ID,
		MerchantID: a.MerchantID,
		UserID:     a.UserID,
		AmountDue:  a.AmountDue,
	}

	return corea, nil
}

// UpdateParams defines the create parameters.
type UpdateParams struct {
	MerchantID uint
	UserID     uint
	AmountDue  uint
}

// UpdateByID updates the given accounting entry by ID.
func (c *Core) UpdateByID(id uint, params *UpdateParams) (*Accounting, error) {
	// Update accounting entry.
	a, err := c.s.Accounting.UpdateByID(id, &accounting.UpdateParams{
		MerchantID: params.MerchantID,
		UserID:     params.UserID,
		AmountDue:  params.AmountDue,
	})
	if err != nil {
		return nil, err
	}

	// Map to core type.
	corea := &Accounting{
		ID:         a.ID,
		MerchantID: a.MerchantID,
		UserID:     a.UserID,
		AmountDue:  a.AmountDue,
	}

	return corea, nil
}
