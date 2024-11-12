package invoice

import (
	"dddstructure/proto"
	serverrors "dddstructure/service/errors"
	"dddstructure/service/interfaces"
	"dddstructure/storage"
	"dddstructure/storage/invoice"
)

// idCounter handles increasing the ID.
var idCounter uint = 1

// Service defines the invoice service.
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

// Create creates a new invoice.
func (s *Service) Create(params *proto.InvoiceCreateParams) (*proto.Invoice, error) {
	// Validate parameters.
	if err := s.ValidateCreateParams(params); err != nil {
		return nil, err
	}

	// Handle ID.
	if params.ID == 0 {
		params.ID = idCounter
		idCounter++
	}

	// Create an invoice.
	storagei, err := s.storage.Invoice.Create(&invoice.Invoice{
		ID:     params.ID,
		UserID: params.UserID,
		BillTo: invoice.BillTo{
			FirstName: params.BillTo.FirstName,
			LastName:  params.BillTo.LastName,
		},
		PayTo: invoice.PayTo{
			FirstName: params.PayTo.FirstName,
			LastName:  params.PayTo.LastName,
		},
		AmountDue:  params.AmountDue,
		AmountPaid: 0,
		Status:     "pending",
	})
	if err != nil {
		return nil, err
	}

	// Map to service type.
	servicei := &proto.Invoice{
		ID:     storagei.ID,
		UserID: storagei.UserID,
		BillTo: proto.InvoiceBillTo{
			FirstName: storagei.BillTo.FirstName,
			LastName:  storagei.BillTo.LastName,
		},
		PayTo: proto.InvoicePayTo{
			FirstName: storagei.PayTo.FirstName,
			LastName:  storagei.PayTo.LastName,
		},
		AmountDue:  storagei.AmountDue,
		AmountPaid: storagei.AmountPaid,
		Status:     storagei.Status,
	}

	return servicei, nil
}

// Get gets a set of invoices.
func (s *Service) Get(params *proto.InvoiceGetParams) ([]*proto.Invoice, error) {
	// Validate parameters.
	if err := s.ValidateGetParams(params); err != nil {
		return nil, err
	}

	// Build the invoice get parameters.
	getParams := &invoice.GetParams{}

	// Check ID.
	if params.ID != nil {
		getParams.ID = params.ID
	}

	// Check user ID.
	if params.UserID != nil {
		getParams.UserID = params.UserID
	}

	// Check status.
	if params.Status != nil {
		getParams.Status = params.Status
	}

	// Get invoices from storage.
	storageis, err := s.storage.Invoice.Get(getParams)
	if err != nil {
		return nil, err
	}

	// Create a new invoices slice.
	invoices := []*proto.Invoice{}

	// Loop through the set of invoices.
	for _, i := range storageis {
		// Create a new invoice.
		invoice := &proto.Invoice{
			ID:     i.ID,
			UserID: i.UserID,
			BillTo: proto.InvoiceBillTo{
				FirstName: i.BillTo.FirstName,
				LastName:  i.BillTo.LastName,
			},
			PayTo: proto.InvoicePayTo{
				FirstName: i.PayTo.FirstName,
				LastName:  i.PayTo.LastName,
			},
			AmountDue:  i.AmountDue,
			AmountPaid: i.AmountPaid,
			Status:     i.Status,
		}

		// Add to invoices slice
		invoices = append(invoices, invoice)
	}

	return invoices, nil
}

// GetCount gets the count of a set of invoices.
func (s *Service) GetCount(params *proto.InvoiceGetParams) (uint, error) {
	// Validate parameters.
	if err := s.ValidateGetParams(params); err != nil {
		return 0, err
	}

	// Build the invoice get parameters.
	getParams := &invoice.GetParams{}

	// Check ID.
	if params.ID != nil {
		getParams.ID = params.ID
	}

	// Check user ID.
	if params.UserID != nil {
		getParams.UserID = params.UserID
	}

	// Check status.
	if params.Status != nil {
		getParams.Status = params.Status
	}

	// Get invoices count from storage.
	count, err := s.storage.Invoice.GetCount(getParams)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetByID gets an invoice by the given ID.
func (s *Service) GetByID(id uint) (*proto.Invoice, error) {
	// Get invoice by ID.
	storagei, err := s.storage.Invoice.GetByID(id)
	if err != nil {
		if err == invoice.ErrInvoiceNotFound {
			return nil, serverrors.ErrInvoiceNotFound
		}

		// Log here.
		return nil, err
	}

	// Map to service type.
	servicei := &proto.Invoice{
		ID:     storagei.ID,
		UserID: storagei.UserID,
		BillTo: proto.InvoiceBillTo{
			FirstName: storagei.BillTo.FirstName,
			LastName:  storagei.BillTo.LastName,
		},
		PayTo: proto.InvoicePayTo{
			FirstName: storagei.PayTo.FirstName,
			LastName:  storagei.PayTo.LastName,
		},
		AmountDue:  storagei.AmountDue,
		AmountPaid: storagei.AmountPaid,
		Status:     storagei.Status,
	}

	return servicei, nil
}

// Update handles updating an invoice.
func (s *Service) Update(params *proto.InvoiceUpdateParams) error {
	// Validate parameters.
	if err := s.ValidateUpdateParams(params); err != nil {
		return err
	}

	// Get invoice from storage.
	storagei, err := s.storage.Invoice.GetByID(*params.ID)
	if err != nil {
		return err
	}

	// Handle amount due.
	if params.AmountDue != nil {
		storagei.AmountDue = *params.AmountDue
	}

	// Handle amount paid.
	if params.AmountPaid != nil {
		storagei.AmountPaid = *params.AmountPaid
	}

	// Handle status.
	if params.Status != nil {
		storagei.Status = *params.Status
	}

	return s.storage.Invoice.Update(storagei)
}

// Pay handles paying an invoice.
func (s *Service) Pay(id uint, params *proto.InvoicePayParams) (*proto.Invoice, error) {
	// Validate parameters.
	if err := s.ValidatePayParams(params); err != nil {
		return nil, err
	}

	// Get the invoice.
	servicei, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Pay the invoice using the transaction service.
	t, err := s.services.Transaction.Process(&proto.TransactionProcessParams{
		UserID:    servicei.UserID,
		Type:      "sale",
		Amount:    params.Amount,
		InvoiceID: id,
	})
	if err != nil {
		return nil, err
	}

	// Update the invoice.
	servicei.AmountPaid = t.AmountCaptured
	servicei.AmountDue -= t.AmountCaptured
	servicei.Status = "paid"

	err = s.Update(&proto.InvoiceUpdateParams{
		ID:         &servicei.ID,
		AmountDue:  &servicei.AmountDue,
		AmountPaid: &servicei.AmountPaid,
		Status:     &servicei.Status,
	})
	if err != nil {
		return nil, err
	}

	return servicei, nil
}
