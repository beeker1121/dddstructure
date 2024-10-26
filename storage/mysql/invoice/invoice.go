package invoice

import (
	"database/sql"
	"fmt"

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

	fmt.Println("Created invoice and added to MySQL database...")

	return inv, nil
}

// GetByID gets an invoice by the given ID.
func (db *Database) GetByID(id uint) (*invoice.Invoice, error) {
	i, ok := invoiceMap[id]
	if !ok {
		return nil, invoice.ErrInvoiceNotFound
	}

	fmt.Println("Got invoice from MySQL database...")

	return i, nil
}

// Update updates an invoice.
func (db *Database) Update(i *invoice.Invoice) error {
	invoiceMap[i.ID] = i

	return nil
}
