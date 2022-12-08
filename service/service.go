package service

import (
	"dddstructure/service/invoice"
	"dddstructure/service/transaction"
	"dddstructure/storage"
)

// Service defines the main business logic service.
type Service struct {
	Invoice     *invoice.Service
	Transaction *transaction.Service
}

// New creates a new service.
func New(s *storage.Storage) *Service {
	return &Service{
		Invoice:     invoice.New(s),
		Transaction: transaction.New(s),
	}
}
