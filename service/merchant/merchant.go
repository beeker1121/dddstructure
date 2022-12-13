package merchant

import (
	"dddstructure/dep"
	"dddstructure/proto"
	"dddstructure/storage"
	"dddstructure/storage/merchant"
)

// idCounter handles increasing the ID.
var idCounter uint = 1

// Service defines the merchant service.
type Service struct {
	s *storage.Storage
}

// New creates a new service.
func New(s *storage.Storage) *Service {
	return &Service{
		s: s,
	}
}

// Create creates a new merchant.
func (s *Service) Create(m *proto.Merchant) (*proto.Merchant, error) {
	// Handle ID.
	if m.ID == 0 {
		m.ID = idCounter
		idCounter++
	}

	// Create a merchant.
	merch, err := s.s.Merchant.Create(&merchant.Merchant{
		ID:    m.ID,
		Name:  m.Name,
		Email: m.Email,
	})
	if err != nil {
		return nil, err
	}

	// Create a user.
	_, err = dep.User.Create(&proto.User{
		AccountTypeID: merch.ID,
		Username:      "johndoe",
	})
	if err != nil {
		return nil, err
	}

	// Map to service type.
	servicem := &proto.Merchant{
		ID:    merch.ID,
		Name:  merch.Name,
		Email: merch.Email,
	}

	return servicem, nil
}

// GetByID gets a merchant by the given ID.
func (s *Service) GetByID(id uint) (*proto.Merchant, error) {
	// Get merchant by ID.
	m, err := s.s.Merchant.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Map to service type.
	servicem := &proto.Merchant{
		ID:    m.ID,
		Name:  m.Name,
		Email: m.Email,
	}

	return servicem, nil
}
