package mysql

import (
	"database/sql"

	"dddstructure/storage"
	"dddstructure/storage/mysql/invoice"
	"dddstructure/storage/mysql/transaction"
	"dddstructure/storage/mysql/user"
)

// New returns a new implementation of storage.Storage that uses MySQL as the
// backend database.
func New(db *sql.DB) *storage.Storage {
	s := &storage.Storage{
		User:        user.New(db),
		Invoice:     invoice.New(db),
		Transaction: transaction.New(db),
	}

	return s
}
