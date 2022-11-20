package mysql

import (
	"database/sql"

	"dddstructure/storage"
	"dddstructure/storage/mysql/merchant"
	"dddstructure/storage/mysql/user"
)

// New returns a new implementation of storage.Storage that uses MySQL as the
// backend database.
func New(db *sql.DB) *storage.Storage {
	s := &storage.Storage{
		Merchant: merchant.New(db),
		User:     user.New(db),
	}

	return s
}
