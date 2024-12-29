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
	model := models.Invoice{
		ID:              i.ID,
		UserID:          i.UserID,
		BillToFirstName: i.BillTo.FirstName,
		BillToLastName:  i.BillTo.LastName,
		PayToFirstName:  i.PayTo.FirstName,
		PayToLastName:   i.PayTo.LastName,
		AmountDue:       i.AmountDue,
		AmountPaid:      i.AmountPaid,
		Status:          models.InvoicesStatus(i.Status),
	}

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
		i := &invoice.Invoice{
			ID:     mi.ID,
			UserID: mi.UserID,
			BillTo: invoice.BillTo{
				FirstName: mi.BillToFirstName,
				LastName:  mi.BillToLastName,
			},
			PayTo: invoice.PayTo{
				FirstName: mi.PayToFirstName,
				LastName:  mi.PayToLastName,
			},
			AmountDue:  mi.AmountDue,
			AmountPaid: mi.AmountPaid,
			Status:     mi.Status.String(),
		}

		invoices = append(invoices, i)
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
	i := &invoice.Invoice{
		ID:     modeli.ID,
		UserID: modeli.UserID,
		BillTo: invoice.BillTo{
			FirstName: modeli.BillToFirstName,
			LastName:  modeli.BillToLastName,
		},
		PayTo: invoice.PayTo{
			FirstName: modeli.PayToFirstName,
			LastName:  modeli.PayToLastName,
		},
		AmountDue:  modeli.AmountDue,
		AmountPaid: modeli.AmountPaid,
		Status:     modeli.Status.String(),
	}

	return i, nil
}

// Update updates an invoice.
func (db *Database) Update(i *invoice.Invoice) error {
	// Map to model.
	model := models.Invoice{
		ID:              i.ID,
		UserID:          i.UserID,
		BillToFirstName: i.BillTo.FirstName,
		BillToLastName:  i.BillTo.LastName,
		PayToFirstName:  i.PayTo.FirstName,
		PayToLastName:   i.PayTo.LastName,
		AmountDue:       i.AmountDue,
		AmountPaid:      i.AmountPaid,
		Status:          models.InvoicesStatus(i.Status),
	}

	// Update in database.
	_, err := model.Update(context.Background(), db.db, boil.Infer())
	if err != nil {
		return err
	}

	return nil
}
