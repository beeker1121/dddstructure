package storage

import (
	"dddstructure/storage/user"
)

// Storage defines the storage system.
type Storage struct {
	User user.Database
}

// New returns a new storage.
func New(storage *Storage) *Storage {
	return storage
}
