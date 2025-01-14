package invoice

import (
	"database/sql"

	"dddstructure/storage/invoice"
)

// invoiceMap acts as a mock MySQL database for invoices.
var invoiceMap map[uint]*invoice.Invoice = make(map[uint]*invoice.Invoice)

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
	inv := &invoice.Invoice{
		ID:         i.ID,
		UserID:     i.UserID,
		BillTo:     i.BillTo,
		PayTo:      i.PayTo,
		AmountDue:  i.AmountDue,
		AmountPaid: i.AmountPaid,
		Status:     i.Status,
	}

	invoiceMap[inv.ID] = inv

	return inv, nil
}

// Get gets a set of invoices.
func (db *Database) Get(params *invoice.GetParams) ([]*invoice.Invoice, error) {
	invoices := []*invoice.Invoice{}
	for _, invoice := range invoiceMap {
		// Handle user ID.
		if params.UserID != nil {
			if invoice.UserID != *params.UserID {
				continue
			}
		}

		invoices = append(invoices, invoice)
	}

	return invoices, nil
}

// GetCount gets the count of a set of invoices.
func (db *Database) GetCount(params *invoice.GetParams) (uint, error) {
	invoices := []*invoice.Invoice{}
	for _, invoice := range invoiceMap {
		// Handle user ID.
		if params.UserID != nil {
			if invoice.UserID != *params.UserID {
				continue
			}
		}

		invoices = append(invoices, invoice)
	}

	return uint(len(invoices)), nil
}

// GetByID gets an invoice by the given ID.
func (db *Database) GetByID(id uint) (*invoice.Invoice, error) {
	i, ok := invoiceMap[id]
	if !ok {
		return nil, invoice.ErrInvoiceNotFound
	}

	return i, nil
}

// GetByPublicHash gets an invoice by the given public hash.
func (db *Database) GetByPublicHash(hash string) (*invoice.Invoice, error) {
	for _, i := range invoiceMap {
		if i.PublicHash == hash {
			return i, nil
		}
	}

	return nil, invoice.ErrInvoiceNotFound
}

// Update updates an invoice.
func (db *Database) Update(i *invoice.Invoice) (*invoice.Invoice, error) {
	invoiceMap[i.ID] = i

	return i, nil
}

// Delete deletes an invoice.
func (db *Database) Delete(id uint) error {
	delete(invoiceMap, id)

	return nil
}
