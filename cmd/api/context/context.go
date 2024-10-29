package context

import (
	"log"

	"dddstructure/cmd/api/config"
	"dddstructure/service"
)

type Context struct {
	Config  *config.Config
	Logger  *log.Logger
	Service *service.Service
}

func New(config *config.Config, logger *log.Logger, services *service.Service) *Context {
	return &Context{
		Config:  config,
		Logger:  logger,
		Service: services,
	}
}
