package storage

import (
	"dddstructure/storage/invoice"
	"dddstructure/storage/transaction"
)

// Storage defines the storage system.
type Storage struct {
	Invoice     invoice.Database
	Transaction transaction.Database
}

// New returns a new storage.
func New(storage *Storage) *Storage {
	return storage
}
