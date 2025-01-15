package invoice

import (
	"log"

	"dddstructure/proto"
	serverrors "dddstructure/service/errors"
	"dddstructure/service/interfaces"
	"dddstructure/storage"
	"dddstructure/storage/invoice"
	"time"

	"github.com/google/uuid"
)

// idCounter handles increasing the ID.
var idCounter uint = 1

// Service defines the invoice service.
type Service struct {
	storage  *storage.Storage
	services *interfaces.Service
	logger   *log.Logger
}

// SetServices sets the services interface.
func (s *Service) SetServices(services *interfaces.Service) {
	s.services = services
}

// New creates a new service.
func New(s *storage.Storage, l *log.Logger) *Service {
	return &Service{
		storage: s,
		logger:  l,
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

	// Handle line items.
	lineItems := []invoice.LineItem{}
	for _, v := range params.LineItems {
		lineItem := invoice.LineItem{
			Name:        v.Name,
			Description: v.Description,
			Quantity:    v.Quantity,
			Price:       v.Price,
		}

		lineItems = append(lineItems, lineItem)
	}

	// Calculate invoice amounts.
	amounts, err := CalculateAmounts(CalculateAmountsParams{
		LineItems: params.LineItems,
		TaxRate:   params.TaxRate,
	})
	if err != nil {
		s.logger.Printf("CalculateAmounts() error: %s\n", err)
		return nil, serverrors.ErrInvoiceCalculatingAmounts
	}

	// Create an invoice.
	paymentMethods := []string{}
	for _, v := range params.PaymentMethods {
		paymentMethods = append(paymentMethods, string(v))
	}

	storagei, err := s.storage.Invoice.Create(&invoice.Invoice{
		ID:            params.ID,
		UserID:        params.UserID,
		PublicHash:    uuid.New().String(),
		InvoiceNumber: params.InvoiceNumber,
		PONumber:      params.PONumber,
		Currency:      params.Currency,
		DueDate:       params.DueDate,
		Message:       params.Message,
		BillTo: invoice.BillTo{
			FirstName:    params.BillTo.FirstName,
			LastName:     params.BillTo.LastName,
			Company:      params.BillTo.Company,
			AddressLine1: params.BillTo.AddressLine1,
			AddressLine2: params.BillTo.AddressLine2,
			City:         params.BillTo.City,
			State:        params.BillTo.State,
			PostalCode:   params.BillTo.PostalCode,
			Country:      params.BillTo.Country,
			Email:        params.BillTo.Email,
			Phone:        params.BillTo.Phone,
		},
		PayTo: invoice.PayTo{
			FirstName:    params.PayTo.FirstName,
			LastName:     params.PayTo.LastName,
			Company:      params.PayTo.Company,
			AddressLine1: params.PayTo.AddressLine1,
			AddressLine2: params.PayTo.AddressLine2,
			City:         params.PayTo.City,
			State:        params.PayTo.State,
			PostalCode:   params.PayTo.PostalCode,
			Country:      params.PayTo.Country,
			Email:        params.PayTo.Email,
			Phone:        params.PayTo.Phone,
		},
		LineItems:      lineItems,
		PaymentMethods: paymentMethods,
		TaxRate:        params.TaxRate,
		AmountDue:      amounts.AmountDue,
		AmountPaid:     0,
		Status:         "pending",
		CreatedAt:      time.Now().UTC(),
	})
	if err != nil {
		s.logger.Printf("storage.Invoice.Create() error: %s\n", err)
		return nil, err
	}

	return storageToProto(storagei), nil
}

// Get gets a set of invoices.
func (s *Service) Get(params *proto.InvoiceGetParams) ([]*proto.Invoice, error) {
	// Validate parameters.
	if err := s.ValidateGetParams(params); err != nil {
		return nil, err
	}

	// Build the invoice get parameters.
	getParams := &invoice.GetParams{
		Offset: params.Offset,
		Limit:  params.Limit,
	}

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

	// Check created at.
	if params.CreatedAt != nil {
		getParams.CreatedAt = &invoice.GetParamsCreatedAt{}
		if params.CreatedAt.StartDate != nil {
			getParams.CreatedAt.StartDate = params.CreatedAt.StartDate
		}
		if params.CreatedAt.EndDate != nil {
			getParams.CreatedAt.EndDate = params.CreatedAt.EndDate
		}
	}

	// Get invoices from storage.
	storageis, err := s.storage.Invoice.Get(getParams)
	if err != nil {
		s.logger.Printf("storage.Invoice.Get() error: %s\n", err)
		return nil, err
	}

	// Create a new invoices slice.
	invoices := []*proto.Invoice{}

	// Loop through the set of invoices.
	for _, i := range storageis {
		// Create a new invoice.
		invoice := storageToProto(i)

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
		s.logger.Printf("storage.Invoice.GetCount() error: %s\n", err)
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

		s.logger.Printf("storage.Invoice.GetByID() error: %s\n", err)
		return nil, err
	}

	return storageToProto(storagei), nil
}

// GetByID gets an invoice by the given ID and user ID.
func (s *Service) GetByIDAndUserID(id, userID uint) (*proto.Invoice, error) {
	// Get invoice by ID.
	storagei, err := s.storage.Invoice.GetByID(id)
	if err != nil {
		if err == invoice.ErrInvoiceNotFound {
			return nil, serverrors.ErrInvoiceNotFound
		}

		s.logger.Printf("storage.Invoice.GetByID() error: %s\n", err)
		return nil, err
	}

	// Check user ID.
	if storagei.UserID != userID {
		return nil, serverrors.ErrInvoiceNotFound
	}

	return storageToProto(storagei), nil
}

// GetByPublicHash gets an invoice by the given public hash.
func (s *Service) GetByPublicHash(hash string) (*proto.Invoice, error) {
	// Get invoice by ID.
	storagei, err := s.storage.Invoice.GetByPublicHash(hash)
	if err != nil {
		if err == invoice.ErrInvoiceNotFound {
			return nil, serverrors.ErrInvoiceNotFound
		}

		s.logger.Printf("storage.Invoice.GetByPublicHash() error: %s\n", err)
		return nil, err
	}

	return storageToProto(storagei), nil
}

// Update handles updating an invoice.
func (s *Service) Update(params *proto.InvoiceUpdateParams) (*proto.Invoice, error) {
	// Validate parameters.
	if err := s.ValidateUpdateParams(params); err != nil {
		return nil, err
	}

	// Get invoice from storage.
	storagei, err := s.storage.Invoice.GetByID(*params.ID)
	if err != nil {
		s.logger.Printf("storage.Invoice.GetByID() error: %s\n", err)
		return nil, err
	}

	// Handle invoice number.
	if params.InvoiceNumber != nil {
		storagei.InvoiceNumber = *params.InvoiceNumber
	}

	// Handle PO number.
	if params.PONumber != nil {
		storagei.PONumber = *params.PONumber
	}

	// Handle currency.
	if params.Currency != nil {
		storagei.Currency = *params.Currency
	}

	// Handle due date.
	if params.DueDate != nil {
		storagei.DueDate = *params.DueDate
	}

	// Handle message.
	if params.Message != nil {
		storagei.Message = *params.Message
	}

	// Handle bill to.
	if params.BillTo != nil {
		if params.BillTo.FirstName != nil {
			storagei.BillTo.FirstName = *params.BillTo.FirstName
		}
		if params.BillTo.LastName != nil {
			storagei.BillTo.LastName = *params.BillTo.LastName
		}
		if params.BillTo.Company != nil {
			storagei.BillTo.Company = *params.BillTo.Company
		}
		if params.BillTo.AddressLine1 != nil {
			storagei.BillTo.AddressLine1 = *params.BillTo.AddressLine1
		}
		if params.BillTo.AddressLine2 != nil {
			storagei.BillTo.AddressLine2 = *params.BillTo.AddressLine2
		}
		if params.BillTo.City != nil {
			storagei.BillTo.City = *params.BillTo.City
		}
		if params.BillTo.State != nil {
			storagei.BillTo.State = *params.BillTo.State
		}
		if params.BillTo.PostalCode != nil {
			storagei.BillTo.PostalCode = *params.BillTo.PostalCode
		}
		if params.BillTo.Country != nil {
			storagei.BillTo.Country = *params.BillTo.Country
		}
		if params.BillTo.Email != nil {
			storagei.BillTo.Email = *params.BillTo.Email
		}
		if params.BillTo.Phone != nil {
			storagei.BillTo.Phone = *params.BillTo.Phone
		}
	}

	// Handle pay to.
	if params.PayTo != nil {
		if params.PayTo.FirstName != nil {
			storagei.PayTo.FirstName = *params.PayTo.FirstName
		}
		if params.PayTo.LastName != nil {
			storagei.PayTo.LastName = *params.PayTo.LastName
		}
		if params.PayTo.Company != nil {
			storagei.PayTo.Company = *params.PayTo.Company
		}
		if params.PayTo.AddressLine1 != nil {
			storagei.PayTo.AddressLine1 = *params.PayTo.AddressLine1
		}
		if params.PayTo.AddressLine2 != nil {
			storagei.PayTo.AddressLine2 = *params.PayTo.AddressLine2
		}
		if params.PayTo.City != nil {
			storagei.PayTo.City = *params.PayTo.City
		}
		if params.PayTo.State != nil {
			storagei.PayTo.State = *params.PayTo.State
		}
		if params.PayTo.PostalCode != nil {
			storagei.PayTo.PostalCode = *params.PayTo.PostalCode
		}
		if params.PayTo.Country != nil {
			storagei.PayTo.Country = *params.PayTo.Country
		}
		if params.PayTo.Email != nil {
			storagei.PayTo.Email = *params.PayTo.Email
		}
		if params.PayTo.Phone != nil {
			storagei.PayTo.Phone = *params.PayTo.Phone
		}
	}

	// Handle line items.
	if params.LineItems != nil {
		lineItems := []invoice.LineItem{}
		for _, v := range *params.LineItems {
			lineItem := invoice.LineItem{
				Name:        v.Name,
				Description: v.Description,
				Quantity:    v.Quantity,
				Price:       v.Price,
			}

			lineItems = append(lineItems, lineItem)
		}

		storagei.LineItems = lineItems
	}

	// Handle payment methods.
	if params.PaymentMethods != nil {
		paymentMethods := []string{}
		for _, v := range *params.PaymentMethods {
			paymentMethods = append(paymentMethods, string(v))
		}

		storagei.PaymentMethods = paymentMethods
	}

	// Handle tax rate.
	if params.TaxRate != nil {
		storagei.TaxRate = *params.TaxRate
	}

	// Calculate invoice amounts.
	amounts, err := CalculateAmounts(CalculateAmountsParams{
		LineItems: storageLineItemsToProto(storagei.LineItems),
		TaxRate:   storagei.TaxRate,
	})
	if err != nil {
		s.logger.Printf("CalculateAmounts() error: %s\n", err)
		return nil, serverrors.ErrInvoiceCalculatingAmounts
	}

	// Set invoice amounts.
	storagei.AmountDue = amounts.AmountDue

	// Update the invoice.
	storagei, err = s.storage.Invoice.Update(storagei)
	if err != nil {
		s.logger.Printf("storage.Invoice.Update() error: %s\n", err)
		return nil, err
	}

	return storageToProto(storagei), nil
}

// UpdateByIDAndUserID handles updating an invoice by given ID and User ID.
func (s *Service) UpdateByIDAndUserID(params *proto.InvoiceUpdateParams) (*proto.Invoice, error) {
	// Get by ID and user ID.
	_, err := s.GetByIDAndUserID(*params.ID, *params.UserID)
	if err != nil {
		return nil, err
	}

	// Call update.
	return s.Update(params)
}

// UpdateForTransaction handles updating an invoice for a transaction.
func (s *Service) UpdateForTransaction(params *proto.InvoiceUpdateForTransactionParams) (*proto.Invoice, error) {
	// Validate parameters.
	if err := s.ValidateUpdateForTransactionParams(params); err != nil {
		return nil, err
	}

	// Get invoice from storage.
	storagei, err := s.storage.Invoice.GetByID(*params.ID)
	if err != nil {
		s.logger.Printf("storage.Invoice.GetByID() error: %s\n", err)
		return nil, err
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

	// Update the invoice.
	storagei, err = s.storage.Invoice.Update(storagei)
	if err != nil {
		s.logger.Printf("storage.Invoice.Update() error: %s\n", err)
		return nil, err
	}

	return storageToProto(storagei), nil
}

// Delete deletes an invoice by the given ID.
func (s *Service) Delete(id uint) error {
	// Delete invoice by ID.
	err := s.storage.Invoice.Delete(id)
	if err != nil {
		if err == invoice.ErrInvoiceNotFound {
			return serverrors.ErrInvoiceNotFound
		}

		s.logger.Printf("storage.Invoice.Delete() error: %s\n", err)
		return err
	}

	return nil
}

// Pay handles paying an invoice.
func (s *Service) Pay(id uint, params *proto.InvoicePayParams) (*proto.Invoice, error) {
	// Validate parameters.
	if err := s.ValidatePayParams(params); err != nil {
		return nil, err
	}

	// Get the invoice.
	storagei, err := s.storage.Invoice.GetByID(id)
	if err != nil {
		s.logger.Printf("storage.Invoice.GetByID() error: %s\n", err)
		return nil, err
	}

	// Check invoice status.
	if storagei.Status != "pending" {
		return nil, serverrors.ErrInvoiceStatusNotPending
	}

	// Pay the invoice using the transaction service.
	t, err := s.services.Transaction.Process(&proto.TransactionProcessParams{
		UserID:    storagei.UserID,
		Type:      "sale",
		Amount:    params.Amount,
		InvoiceID: id,
	})
	if err != nil {
		return nil, err
	}

	// Update the invoice.
	storagei.AmountPaid = t.AmountCaptured
	storagei.AmountDue -= t.AmountCaptured
	storagei.Status = "paid"

	storagei, err = s.storage.Invoice.Update(storagei)
	if err != nil {
		s.logger.Printf("storage.Invoice.Update() error: %s\n", err)
		return nil, err
	}

	return storageToProto(storagei), nil
}

// storageLineItemsToProto handles mappings the storage invoice line items type
// to the proto invoice line items type.
func storageLineItemsToProto(li []invoice.LineItem) []proto.InvoiceLineItem {
	lineItems := []proto.InvoiceLineItem{}
	for _, v := range li {
		lineItem := proto.InvoiceLineItem{
			Name:        v.Name,
			Description: v.Description,
			Quantity:    v.Quantity,
			Price:       v.Price,
		}

		lineItems = append(lineItems, lineItem)
	}

	return lineItems
}

// storageToProto handles mapping a storage invoice type to the proto invoice
// type.
func storageToProto(s *invoice.Invoice) *proto.Invoice {
	// Handle line items.
	lineItems := storageLineItemsToProto(s.LineItems)

	paymentMethods := []proto.InvoicePaymentMethod{}
	for _, v := range s.PaymentMethods {
		paymentMethods = append(paymentMethods, proto.InvoicePaymentMethod(v))
	}

	return &proto.Invoice{
		ID:            s.ID,
		UserID:        s.UserID,
		PublicHash:    s.PublicHash,
		InvoiceNumber: s.InvoiceNumber,
		PONumber:      s.PONumber,
		Currency:      s.Currency,
		DueDate:       s.DueDate,
		Message:       s.Message,
		BillTo: proto.InvoiceBillTo{
			FirstName:    s.BillTo.FirstName,
			LastName:     s.BillTo.LastName,
			Company:      s.BillTo.Company,
			AddressLine1: s.BillTo.AddressLine1,
			AddressLine2: s.BillTo.AddressLine2,
			City:         s.BillTo.City,
			State:        s.BillTo.State,
			PostalCode:   s.BillTo.PostalCode,
			Country:      s.BillTo.Country,
			Email:        s.BillTo.Email,
			Phone:        s.BillTo.Phone,
		},
		PayTo: proto.InvoicePayTo{
			FirstName:    s.PayTo.FirstName,
			LastName:     s.PayTo.LastName,
			Company:      s.PayTo.Company,
			AddressLine1: s.PayTo.AddressLine1,
			AddressLine2: s.PayTo.AddressLine2,
			City:         s.PayTo.City,
			State:        s.PayTo.State,
			PostalCode:   s.PayTo.PostalCode,
			Country:      s.PayTo.Country,
			Email:        s.PayTo.Email,
			Phone:        s.PayTo.Phone,
		},
		LineItems:      lineItems,
		PaymentMethods: paymentMethods,
		TaxRate:        s.TaxRate,
		AmountDue:      s.AmountDue,
		AmountPaid:     s.AmountPaid,
		Status:         s.Status,
		CreatedAt:      s.CreatedAt,
	}
}
