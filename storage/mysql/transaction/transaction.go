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
func (db *Database) Create(params *transaction.CreateParams) (*transaction.Transaction, error) {
	t := &transaction.Transaction{
		ID:             params.ID,
		MerchantID:     params.MerchantID,
		Type:           params.Type,
		CardType:       params.CardType,
		AmountCaptured: params.AmountCaptured,
	}

	transactionMap[t.ID] = t

	fmt.Println("Created transaction and added to MySQL database...")

	return t, nil
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
