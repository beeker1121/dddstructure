package context

import (
	"log"

	"dddstructure/cmd/api/config"
	"dddstructure/service"
)

// Context defines the API context.
type Context struct {
	Config  *config.Config
	Logger  *log.Logger
	Service *service.Service
}

// New returns a new API context.
func New(config *config.Config, logger *log.Logger, services *service.Service) *Context {
	return &Context{
		Config:  config,
		Logger:  logger,
		Service: services,
	}
}
