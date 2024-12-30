package invoice

import (
	"context"
	"database/sql"

	"dddstructure/storage/invoice"
	"dddstructure/storage/mysql/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// Database defines the database.
type Database struct {
	db *sql.DB
}

// New creates a new database.
func New(db *sql.DB) *Database {
	return &Database{
		db: db,
	}
}

// Create creates a new invoice.
func (db *Database) Create(i *invoice.Invoice) (*invoice.Invoice, error) {
	// Map to model.
	model := storageToModel(i)

	// Insert into database.
	err := model.Insert(context.Background(), db.db, boil.Infer())
	if err != nil {
		return nil, err
	}

	return i, nil
}

// Get gets a set of invoices.
func (db *Database) Get(params *invoice.GetParams) ([]*invoice.Invoice, error) {
	var filter []qm.QueryMod

	// Handle get params.
	if params.UserID != nil {
		filter = append(filter, qm.Where("user_id=?", params.UserID))
	}

	// Get from database.
	modelInvoices, err := models.Invoices(filter...).All(context.Background(), db.db)
	if err != nil {
		return nil, err
	}

	// Build invoices slice.
	invoices := []*invoice.Invoice{}
	for _, mi := range modelInvoices {
		i := modelToStorage(mi)

		invoices = append(invoices, &i)
	}

	return invoices, nil
}

// GetCount gets the count of a set of invoices.
func (db *Database) GetCount(params *invoice.GetParams) (uint, error) {
	var filter []qm.QueryMod

	// Handle get params.
	if params.UserID != nil {
		filter = append(filter, qm.Where("user_id=?", params.UserID))
	}

	// Get from database.
	count, err := models.Invoices(filter...).Count(context.Background(), db.db)
	if err != nil {
		return 0, err
	}

	return uint(count), nil
}

// GetByID gets an invoice by the given ID.
func (db *Database) GetByID(id uint) (*invoice.Invoice, error) {
	modeli, err := models.Invoices(qm.Where("id=?", id)).One(context.Background(), db.db)
	if err == sql.ErrNoRows {
		return nil, invoice.ErrInvoiceNotFound
	} else if err != nil {
		return nil, err
	}

	// Map to invoice type.
	i := modelToStorage(modeli)

	return &i, nil
}

// Update updates an invoice.
func (db *Database) Update(i *invoice.Invoice) (*invoice.Invoice, error) {
	// Map to model.
	model := storageToModel(i)

	// Update in database.
	_, err := model.Update(context.Background(), db.db, boil.Infer())
	if err != nil {
		return nil, err
	}

	return i, nil
}

// storageToModel handles mapping a storage invoice type to the model invoice
// type.
func storageToModel(i *invoice.Invoice) models.Invoice {
	return models.Invoice{
		ID:                 i.ID,
		UserID:             i.UserID,
		InvoiceNumber:      i.InvoiceNumber,
		PoNumber:           i.PONumber,
		Currency:           i.Currency,
		DueDate:            i.DueDate,
		Message:            i.Message,
		BillToFirstName:    i.BillTo.FirstName,
		BillToLastName:     i.BillTo.LastName,
		BillToCompany:      i.BillTo.Company,
		BillToAddressLine1: i.BillTo.AddressLine1,
		BillToAddressLine2: i.BillTo.AddressLine2,
		BillToCity:         i.BillTo.City,
		BillToState:        i.BillTo.State,
		BillToPostalCode:   i.BillTo.PostalCode,
		BillToCountry:      i.BillTo.Country,
		BillToEmail:        i.BillTo.Email,
		BillToPhone:        i.BillTo.Phone,
		PayToFirstName:     i.PayTo.FirstName,
		PayToLastName:      i.PayTo.LastName,
		PayToCompany:       i.PayTo.Company,
		PayToAddressLine1:  i.PayTo.AddressLine1,
		PayToAddressLine2:  i.PayTo.AddressLine2,
		PayToCity:          i.PayTo.City,
		PayToState:         i.PayTo.State,
		PayToPostalCode:    i.PayTo.PostalCode,
		PayToCountry:       i.PayTo.Country,
		PayToEmail:         i.PayTo.Email,
		PayToPhone:         i.PayTo.Phone,
		TaxRate:            i.TaxRate,
		AmountDue:          i.AmountDue,
		AmountPaid:         i.AmountPaid,
		Status:             models.InvoicesStatus(i.Status),
	}
}

// modelToStorage handles mapping a model invoice type to the storage invoice
// type.
func modelToStorage(i *models.Invoice) invoice.Invoice {
	return invoice.Invoice{
		ID:            i.ID,
		UserID:        i.UserID,
		InvoiceNumber: i.InvoiceNumber,
		PONumber:      i.PoNumber,
		Currency:      i.Currency,
		DueDate:       i.DueDate,
		Message:       i.Message,
		BillTo: invoice.BillTo{
			FirstName:    i.BillToFirstName,
			LastName:     i.BillToLastName,
			Company:      i.BillToCompany,
			AddressLine1: i.BillToAddressLine1,
			AddressLine2: i.BillToAddressLine2,
			City:         i.BillToCity,
			State:        i.BillToState,
			PostalCode:   i.BillToPostalCode,
			Country:      i.BillToCountry,
			Email:        i.BillToEmail,
			Phone:        i.BillToPhone,
		},
		PayTo: invoice.PayTo{
			FirstName:    i.PayToFirstName,
			LastName:     i.PayToLastName,
			Company:      i.PayToCompany,
			AddressLine1: i.PayToAddressLine1,
			AddressLine2: i.PayToAddressLine2,
			City:         i.PayToCity,
			State:        i.PayToState,
			PostalCode:   i.PayToPostalCode,
			Country:      i.PayToCountry,
			Email:        i.PayToEmail,
			Phone:        i.PayToPhone,
		},
		TaxRate:    i.TaxRate,
		AmountDue:  i.AmountDue,
		AmountPaid: i.AmountPaid,
		Status:     i.Status.String(),
	}
}
