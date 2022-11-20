package storage

import (
	"dddstructure/storage/merchant"
	"dddstructure/storage/user"
)

// Storage defines the storage system.
type Storage struct {
	Merchant merchant.Database
	User     user.Database
}

// New returns a new storage.
func New(storage *Storage) *Storage {
	return storage
}
