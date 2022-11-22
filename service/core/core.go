package core

import (
	"dddstructure/service/core/accounting"
	"dddstructure/service/core/billing"
	"dddstructure/service/core/merchant"
	"dddstructure/service/core/user"
	"dddstructure/storage"
)

// Core defines the core service.
type Core struct {
	Merchant   *merchant.Core
	User       *user.Core
	Accounting *accounting.Core
	Billing    *billing.Core
}

// New creates a new core.
func New(s *storage.Storage) *Core {
	return &Core{
		Merchant:   merchant.New(s),
		User:       user.New(s),
		Accounting: accounting.New(s),
		Billing:    billing.New(s),
	}
}
