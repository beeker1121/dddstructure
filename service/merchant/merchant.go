package merchant

import (
	"dddstructure/service/core"
	"dddstructure/service/core/merchant"
	"dddstructure/service/core/user"
)

// Service defines the merchant service.
type Service struct {
	c *core.Core
}

// New creates a new service.
func New(c *core.Core) *Service {
	return &Service{
		c: c,
	}
}

// Merchant defines a merchant.
type Merchant struct {
	ID    uint
	Name  string
	Email string
}

// CreateParams defines the Create parameters.
type CreateParams struct {
	ID    uint
	Name  string
	Email string
}

// Create creates a new merchant.
//
// This showcases using the core services, ie. core merchant create and core
// user create, to handle overall business logic - when we create a new
// merchant, we also need to create a new user for that merchant.
func (s *Service) Create(params *CreateParams) (*Merchant, error) {
	// Create a merchant.
	m, err := s.c.Merchant.Create(&merchant.CreateParams{
		ID:    params.ID,
		Name:  params.Name,
		Email: params.Email,
	})
	if err != nil {
		return nil, err
	}

	// Create a user for this merchant.
	_, err = s.c.User.Create(&user.CreateParams{
		ID:    params.ID,
		Name:  params.Name,
		Email: params.Email,
	})

	// Map to service type.
	servicem := &Merchant{
		ID:    m.ID,
		Name:  m.Name,
		Email: m.Email,
	}

	return servicem, nil
}

// GetByID gets a merchant by the given ID.
func (s *Service) GetByID(id uint) (*Merchant, error) {
	// Get merchant by ID.
	m, err := s.c.Merchant.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Map to service type.
	servicem := &Merchant{
		ID:    m.ID,
		Name:  m.Name,
		Email: m.Email,
	}

	return servicem, nil
}
