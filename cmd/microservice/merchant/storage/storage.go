package storage

import (
	"dddstructure/storage/merchant"
)

// Storage defines the storage system.
type Storage struct {
	Merchant merchant.Database
}

// New returns a new storage.
func New(storage *Storage) *Storage {
	return storage
}
