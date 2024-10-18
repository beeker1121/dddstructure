package service

import (
	"dddstructure/service/interfaces"
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

func (s *Service) SetServices(services *interfaces.Service) {
	s.Merchant.SetServices(services)
	s.User.SetServices(services)
	s.Invoice.SetServices(services)
	s.Transaction.SetServices(services)
}

// New creates a new service.
func New(s *storage.Storage) *Service {
	// Create services.
	serv := &Service{
		Merchant:    merchant.New(s),
		User:        user.New(s),
		Invoice:     invoice.New(s),
		Transaction: transaction.New(s),
	}

	// Create services interface.
	servi := interfaces.NewService(interfaces.NewServiceParams{
		Merchant:    serv.Merchant,
		User:        serv.User,
		Invoice:     serv.Invoice,
		Transaction: serv.Transaction,
	})

	// Set services interfaces for all services.
	serv.SetServices(servi)

	return serv
}
