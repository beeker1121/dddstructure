package transaction

import (
	"database/sql"
	"errors"
	"fmt"

	"dddstructure/storage/transaction"
)

// transactionMap acts as a mock MySQL database for transactions.
var transactionMap map[uint]*transaction.Transaction = make(map[uint]*transaction.Transaction)

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

// Create creates a new transaction.
func (db *Database) Create(t *transaction.Transaction) (*transaction.Transaction, error) {
	trans := &transaction.Transaction{
		ID:             t.ID,
		UserID:         t.UserID,
		Type:           t.Type,
		CardType:       t.CardType,
		AmountCaptured: t.AmountCaptured,
		InvoiceID:      t.InvoiceID,
		Status:         t.Status,
	}

	transactionMap[trans.ID] = trans

	fmt.Println("Created transaction and added to MySQL database...")

	return trans, nil
}

// GetByID gets a transaction by the given ID.
func (db *Database) GetByID(id uint) (*transaction.Transaction, error) {
	m, ok := transactionMap[id]
	if !ok {
		return nil, errors.New("could not find transaction")
	}

	fmt.Println("Got transaction from MySQL database...")

	return m, nil
}
