package storage

import (
	"dddstructure/storage/invoice"
	"dddstructure/storage/merchant"
	"dddstructure/storage/transaction"
	"dddstructure/storage/user"
)

// Storage defines the storage system.
type Storage struct {
	Merchant    merchant.Database
	User        user.Database
	Invoice     invoice.Database
	Transaction transaction.Database
}

// New returns a new storage.
func New(storage *Storage) *Storage {
	return storage
}
