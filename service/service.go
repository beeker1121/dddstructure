package service

import (
	"log/slog"

	"dddstructure/service/interfaces"
	"dddstructure/service/invoice"
	"dddstructure/service/transaction"
	"dddstructure/service/user"
	"dddstructure/storage"
)

// Service defines the main business logic service.
type Service struct {
	User        *user.Service
	Invoice     *invoice.Service
	Transaction *transaction.Service
}

// SetServices sets the services interface for all individual services.
//
// This is done so each individual service has access to all other top level
// services in the app. One service will be able to call the function of
// another service and vice versa, and this method gets around cyclical imports
// within Go.
func (s *Service) SetServices(services *interfaces.Service) {
	s.User.SetServices(services)
	s.Invoice.SetServices(services)
	s.Transaction.SetServices(services)
}

// New creates a new service.
func New(s *storage.Storage, l *slog.Logger) *Service {
	// Create services.
	serv := &Service{
		User:        user.New(s, l),
		Invoice:     invoice.New(s, l),
		Transaction: transaction.New(s, l),
	}

	// Create services interface.
	servi := interfaces.NewService(interfaces.NewServiceParams{
		User:        serv.User,
		Invoice:     serv.Invoice,
		Transaction: serv.Transaction,
	})

	// Set services interfaces for all services.
	serv.SetServices(servi)

	return serv
}
