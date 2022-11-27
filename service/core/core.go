package core

import (
	"dddstructure/service/core/accounting"
	"dddstructure/service/core/billing"
	"dddstructure/service/core/invoice"
	"dddstructure/service/core/merchant"
	"dddstructure/service/core/transaction"
	"dddstructure/service/core/user"
	"dddstructure/storage"
)

// Core defines the core service.
type Core struct {
	Merchant    *merchant.Core
	User        *user.Core
	Invoice     *invoice.Core
	Transaction *transaction.Core
	Accounting  *accounting.Core
	Billing     *billing.Core
}

// New creates a new core.
func New(s *storage.Storage) *Core {
	return &Core{
		Merchant:    merchant.New(s),
		User:        user.New(s),
		Invoice:     invoice.New(s),
		Transaction: transaction.New(s),
		Accounting:  accounting.New(s),
		Billing:     billing.New(s),
	}
}
