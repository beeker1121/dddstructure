package service

import (
	"dddstructure/service/core"
	"dddstructure/service/merchant"
	"dddstructure/service/user"
)

// Service defines the main business logic service.
type Service struct {
	Merchant *merchant.Service
	User     *user.Service
}

// New creates a new service.
func New(c *core.Core) *Service {
	return &Service{
		Merchant: merchant.New(c),
		User:     user.New(c),
	}
}
