package service

import (
	"dddstructure/service/invoice"
	"dddstructure/service/merchant"
	"dddstructure/service/transaction"
	"dddstructure/service/user"
	"dddstructure/storage"
)

// Service defines the main business logic service.
type Service struct {
	Merchant    *merchant.Service
	User        *user.Service
	Invoice     *invoice.Service
	Transaction *transaction.Service
}

// New creates a new service.
func New(s *storage.Storage) *Service {
	return &Service{
		Merchant:    merchant.New(s),
		User:        user.New(s),
		Invoice:     invoice.New(s),
		Transaction: transaction.New(s),
	}
}
