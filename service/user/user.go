package user

import (
	"dddstructure/proto"
	serverrors "dddstructure/service/errors"
	"dddstructure/service/interfaces"
	"dddstructure/storage"
	"dddstructure/storage/user"
)

// idCounter handles increasing the ID.
var idCounter uint = 1

// Service defines the user service.
type Service struct {
	storage  *storage.Storage
	services *interfaces.Service
}

// SetServices sets the services interface.
func (s *Service) SetServices(services *interfaces.Service) {
	s.services = services
}

// New creates a new service.
func New(s *storage.Storage) *Service {
	return &Service{
		storage: s,
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
	use, err := s.storage.User.Create(&user.User{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	})
	if err != nil {
		return nil, err
	}

	// Map to service type.
	serviceu := &proto.User{
		ID:       use.ID,
		Username: use.Username,
		Email:    u.Email,
	}

	return serviceu, nil
}

// GetByID gets a user by the given ID.
func (s *Service) GetByID(id uint) (*proto.User, error) {
	// Get user by ID.
	u, err := s.storage.User.GetByID(id)
	if err != nil {
		if err == user.ErrUserNotFound {
			return nil, serverrors.ErrUserNotFound
		}

		// Log here.
		return nil, err
	}

	// Map to service type.
	serviceu := &proto.User{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}

	return serviceu, nil
}
