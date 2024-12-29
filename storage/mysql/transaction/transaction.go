package transaction

import (
	"context"
	"database/sql"

	"dddstructure/storage/mysql/models"
	"dddstructure/storage/transaction"

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

// Create creates a new transaction.
func (db *Database) Create(t *transaction.Transaction) (*transaction.Transaction, error) {
	// Map to model.
	model := models.Transaction{
		ID:             t.ID,
		UserID:         t.UserID,
		Type:           models.TransactionsType(t.Type),
		CardType:       t.CardType,
		AmountCaptured: t.AmountCaptured,
		InvoiceID:      t.InvoiceID,
		Status:         models.TransactionsStatus(t.Status),
	}

	// Insert into database.
	err := model.Insert(context.Background(), db.db, boil.Infer())
	if err != nil {
		return nil, err
	}

	return t, nil
}

// GetByID gets a transaction by the given ID.
func (db *Database) GetByID(id uint) (*transaction.Transaction, error) {
	modelt, err := models.Transactions(qm.Where("id=?", id)).One(context.Background(), db.db)
	if err == sql.ErrNoRows {
		return nil, transaction.ErrTransactionNotFound
	} else if err != nil {
		return nil, err
	}

	// Map to transaction type.
	t := &transaction.Transaction{
		ID:             modelt.ID,
		UserID:         modelt.UserID,
		Type:           modelt.Type.String(),
		CardType:       modelt.CardType,
		AmountCaptured: modelt.AmountCaptured,
		InvoiceID:      modelt.InvoiceID,
		Status:         modelt.Status.String(),
	}

	return t, nil
}
