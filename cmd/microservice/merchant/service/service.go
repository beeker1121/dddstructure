package service

import (
	"dddstructure/cmd/microservice/merchant/service/merchant"
	"dddstructure/cmd/microservice/merchant/service/user"
	"dddstructure/cmd/microservice/merchant/storage"
)

// Service defines the main business logic service.
type Service struct {
	Merchant *merchant.Service
	User     *user.Service
}

// New creates a new service.
func New(s *storage.Storage) *Service {
	return &Service{
		Merchant: merchant.New(s),
		User:     user.New(),
	}
}
