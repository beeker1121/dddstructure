package user

import (
	"dddstructure/proto"
	"dddstructure/storage"
	"dddstructure/storage/user"
)

// idCounter handles increasing the ID.
var idCounter uint = 1

// Service defines the user service.
type Service struct {
	s *storage.Storage
}

// New creates a new service.
func New(s *storage.Storage) *Service {
	return &Service{
		s: s,
	}
}

// Create creates a new user.
func (s *Service) Create(u *proto.User) (*proto.User, error) {
	// Handle ID.
	if u.ID == 0 {
		u.ID = idCounter
		idCounter++
	}

	// Create a user.
	use, err := s.s.User.Create(&user.User{
		ID:            u.ID,
		AccountTypeID: u.AccountTypeID,
		Username:      u.Username,
	})
	if err != nil {
		return nil, err
	}

	// Map to service type.
	serviceu := &proto.User{
		ID:       use.ID,
		Username: use.Username,
	}

	return serviceu, nil
}

// GetByID gets a user by the given ID.
func (s *Service) GetByID(id uint) (*proto.User, error) {
	// Get user by ID.
	u, err := s.s.User.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Map to service type.
	serviceu := &proto.User{
		ID:            u.ID,
		AccountTypeID: u.AccountTypeID,
		Username:      u.Username,
	}

	return serviceu, nil
}
