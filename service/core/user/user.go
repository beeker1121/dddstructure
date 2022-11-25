package user

import (
	"dddstructure/storage"
	"dddstructure/storage/user"
)

// idIncrementer handles incrementing the user IDs.
var idIncrementer uint = 1

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
	// Handle creating an ID if one is not present.
	if params.ID == 0 {
		params.ID = idIncrementer
		idIncrementer++
	}

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
		Name:  u.Name,
		Email: u.Email,
	}

	return coreu, nil
}

// GetByID gets a user by the given ID.
func (c *Core) GetByID(id uint) (*User, error) {
	// Get user by ID.
	u, err := c.s.User.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Map to core type.
	coreu := &User{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}

	return coreu, nil
}
