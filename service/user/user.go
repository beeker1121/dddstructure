package user

import (
	"dddstructure/proto"
	serverrors "dddstructure/service/errors"
	"dddstructure/service/interfaces"
	"dddstructure/storage"
	"dddstructure/storage/user"

	"golang.org/x/crypto/bcrypt"
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
func (s *Service) Create(params *proto.UserCreateParams) (*proto.User, error) {
	// Validate parameters.
	if err := s.ValidateCreateParams(params); err != nil {
		return nil, err
	}

	// Handle email.
	pes := serverrors.NewParamErrors()
	_, err := s.storage.User.GetByEmail(params.Email)
	if err == nil {
		pes.Add(serverrors.NewParamError("email", serverrors.ErrUserEmailExists))
	} else if err != nil && err != user.ErrUserNotFound {
		return nil, err
	}

	// Return if there were parameter errors.
	if pes.Length() > 0 {
		return nil, pes
	}

	// Handle ID.
	if params.ID == 0 {
		params.ID = idCounter
		idCounter++
	}

	// Hash the password.
	pwHash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create a user.
	storageu, err := s.storage.User.Create(&user.User{
		ID:       params.ID,
		Email:    params.Email,
		Password: string(pwHash),
	})
	if err != nil {
		return nil, err
	}

	// Map to service type.
	serviceu := &proto.User{
		ID:       storageu.ID,
		Email:    storageu.Email,
		Password: storageu.Password,
	}

	return serviceu, nil
}

// Login checks if a user exists in the database and can log in.
func (s *Service) Login(params *proto.UserLoginParams) (*proto.User, error) {
	// Validate parameters.
	if err := s.ValidateLoginParams(params); err != nil {
		return nil, err
	}

	// Try to pull this user from the database by email.
	storageu, err := s.storage.User.GetByEmail(params.Email)
	if err == user.ErrUserNotFound {
		return nil, serverrors.ErrUserInvalidLogin
	} else if err != nil {
		return nil, err
	}

	// Validate the password.
	if err := bcrypt.CompareHashAndPassword([]byte(storageu.Password), []byte(params.Password)); err != nil {
		return nil, serverrors.ErrUserInvalidLogin
	}

	// Map to service type.
	serviceu := &proto.User{
		ID:       storageu.ID,
		Email:    storageu.Email,
		Password: storageu.Password,
	}

	return serviceu, nil
}

// GetByID gets a user by the given ID.
func (s *Service) GetByID(id uint) (*proto.User, error) {
	// Get user by ID.
	storageu, err := s.storage.User.GetByID(id)
	if err != nil {
		if err == user.ErrUserNotFound {
			return nil, serverrors.ErrUserNotFound
		}

		// Log here.
		return nil, err
	}

	// Map to service type.
	serviceu := &proto.User{
		ID:       storageu.ID,
		Email:    storageu.Email,
		Password: storageu.Password,
	}

	return serviceu, nil
}

// Update handles updating a user.
func (s *Service) Update(params *proto.UserUpdateParams) (*proto.User, error) {
	// Validate parameters.
	if err := s.ValidateUpdateParams(params); err != nil {
		return nil, err
	}

	// Get the user.
	serviceu, err := s.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	// Check email.
	pes := serverrors.NewParamErrors()
	if params.Email != nil && *params.Email != serviceu.Email {
		_, err := s.storage.User.GetByEmail(*params.Email)
		if err == nil {
			pes.Add(serverrors.NewParamError("email", serverrors.ErrUserEmailExists))
		} else if err != nil && err != user.ErrUserNotFound {
			return nil, err
		}
	}

	// Return if there were parameter errors.
	if pes.Length() > 0 {
		return nil, pes
	}

	// Get user from storage.
	storageu, err := s.storage.User.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	// Handle email.
	if params.Email != nil && *params.Email != serviceu.Email {
		storageu.Email = *params.Email
	}

	// Hash the password.
	if params.Password != nil {
		pwHash, err := bcrypt.GenerateFromPassword([]byte(*params.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		storageu.Password = string(pwHash)
	}

	// Update the user.
	storageu, err = s.storage.User.Update(storageu)
	if err != nil {
		return nil, err
	}

	// Map to service type.
	serviceu = &proto.User{
		ID:       storageu.ID,
		Email:    storageu.Email,
		Password: storageu.Password,
	}

	return serviceu, nil
}
