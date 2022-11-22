package accounting

import (
	"database/sql"
	"errors"
	"fmt"

	"dddstructure/storage/accounting"
)

// accountingMap acts as a mock MySQL database for accounting.
var accountingMap map[uint]*accounting.Accounting = make(map[uint]*accounting.Accounting)

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

// Create creates a new accounting entry.
func (db *Database) Create(params *accounting.CreateParams) (*accounting.Accounting, error) {
	a := &accounting.Accounting{
		ID:         params.ID,
		MerchantID: params.MerchantID,
		UserID:     params.UserID,
		AmountDue:  params.AmountDue,
	}

	accountingMap[a.ID] = a

	fmt.Println("Created accounting entry and added to MySQL database...")

	return a, nil
}

// GetByID gets a merchant by the given ID.
func (db *Database) GetByID(id uint) (*accounting.Accounting, error) {
	a, ok := accountingMap[id]
	if !ok {
		return nil, errors.New("could not find accounting entry")
	}

	fmt.Println("Got accounting entry from MySQL database...")

	return a, nil
}
