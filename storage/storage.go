package storage

import (
	"dddstructure/storage/accounting"
	"dddstructure/storage/billing"
	"dddstructure/storage/merchant"
	"dddstructure/storage/user"
)

// Storage defines the storage system.
type Storage struct {
	Merchant   merchant.Database
	User       user.Database
	Accounting accounting.Database
	Billing    billing.Database
}

// New returns a new storage.
func New(storage *Storage) *Storage {
	return storage
}
