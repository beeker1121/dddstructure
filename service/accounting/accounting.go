package accounting

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

// Create creates a new user.
func (s *Service) Create(params *CreateParams) (*Accounting, error) {
	// Create a new accounting entry.
	a, err := s.c.Accounting.Create(&accounting.CreateParams{
		ID:         params.ID,
		MerchantID: params.MerchantID,
		UserID:     params.UserID,
		AmountDue:  params.AmountDue,
	})
	if err != nil {
		return nil, err
	}

	// Map to service type.
	servicea := &Accounting{
		ID:         a.ID,
		MerchantID: a.MerchantID,
		UserID:     a.UserID,
		AmountDue:  a.AmountDue,
	}

	return servicea, nil
}

// GetByID gets an accounting entry by the given ID.
func (s *Service) GetByID(id uint) (*Accounting, error) {
	// Get user by ID.
	a, err := s.c.Accounting.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Map to service type.
	servicea := &Accounting{
		ID:         a.ID,
		MerchantID: a.MerchantID,
		UserID:     a.UserID,
		AmountDue:  a.AmountDue,
	}

	return servicea, nil
}
