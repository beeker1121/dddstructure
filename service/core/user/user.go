package user

import (
	"dddstructure/storage"
	"dddstructure/storage/user"
)

// Core defines the core user service.
type Core struct {
	s *storage.Storage
}

// New creates a new core.
func New(s *storage.Storage) *Core {
	return &Core{
		s: s,
	}
}

// User defines a user.
type User struct {
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

// Create creates a new user.
func (c *Core) Create(params *CreateParams) (*User, error) {
	// Get user by ID.
	u, err := c.s.User.Create(&user.CreateParams{
		ID:    params.ID,
		Name:  params.Name,
		Email: params.Email,
	})
	if err != nil {
		return nil, err
	}

	// Map to core type.
	coreu := &User{
		ID:    u.ID,
		Name:  u.Email,
		Email: u.Email,
	}

	return coreu, nil
}

// GetByID gets a user by the given ID.
func (c *Core) GetByID(id uint) (*User, error) {
	// Get user by ID.
	m, err := c.s.User.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Map to core type.
	corem := &User{
		ID:    m.ID,
		Name:  m.Name,
		Email: m.Email,
	}

	return corem, nil
}
