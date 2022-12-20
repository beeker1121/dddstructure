package service

import (
	"dddstructure/cmd/microservice/user/service/user"
	"dddstructure/cmd/microservice/user/storage"
)

// Service defines the main business logic service.
type Service struct {
	User *user.Service
}

// New creates a new service.
func New(s *storage.Storage) *Service {
	return &Service{
		User: user.New(s),
	}
}
