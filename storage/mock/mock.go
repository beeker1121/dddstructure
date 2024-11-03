package mock

import (
	"database/sql"

	"dddstructure/storage"
	"dddstructure/storage/mock/invoice"
	"dddstructure/storage/mock/transaction"
	"dddstructure/storage/mock/user"
)

// New returns a new implementation of storage.Storage that uses a mock as the
// backend database.
func New(db *sql.DB) *storage.Storage {
	s := &storage.Storage{
		User:        user.New(db),
		Invoice:     invoice.New(db),
		Transaction: transaction.New(db),
	}

	return s
}
