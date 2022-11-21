package merchant

import (
	"dddstructure/storage"
	"dddstructure/storage/merchant"
)

// Core defines the core merchant service.
type Core struct {
	s *storage.Storage
}

// New creates a new core.
func New(s *storage.Storage) *Core {
	return &Core{
		s: s,
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
func (c *Core) Create(params *CreateParams) (*Merchant, error) {
	// Create a merchant.
	m, err := c.s.Merchant.Create(&merchant.CreateParams{
		ID:    params.ID,
		Name:  params.Name,
		Email: params.Email,
	})
	if err != nil {
		return nil, err
	}

	// Map to core type.
	corem := &Merchant{
		ID:    m.ID,
		Name:  m.Name,
		Email: m.Email,
	}

	return corem, nil
}

// GetByID gets a merchant by the given ID.
func (c *Core) GetByID(id uint) (*Merchant, error) {
	// Get merchant by ID.
	m, err := c.s.Merchant.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Map to core type.
	corem := &Merchant{
		ID:    m.ID,
		Name:  m.Name,
		Email: m.Email,
	}

	return corem, nil
}
