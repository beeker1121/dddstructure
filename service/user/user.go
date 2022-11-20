package user

import (
	"dddstructure/service/core"
	"dddstructure/service/core/user"
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
func (s *Service) Create(params *CreateParams) (*User, error) {
	// Create a user.
	u, err := s.c.User.Create(&user.CreateParams{
		ID:    params.ID,
		Name:  params.Name,
		Email: params.Email,
	})
	if err != nil {
		return nil, err
	}

	// Map to service type.
	serviceu := &User{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}

	return serviceu, nil
}

// GetByID gets a user by the given ID.
func (s *Service) GetByID(id uint) (*User, error) {
	// Get user by ID.
	m, err := s.c.User.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Map to service type.
	servicem := &User{
		ID:    m.ID,
		Name:  m.Name,
		Email: m.Email,
	}

	return servicem, nil
}
