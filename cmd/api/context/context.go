package context

import "dddstructure/service"

// Context defines the API context, which acts as a container for all assets
// used by the API.
type Context struct {
	Services *service.Service
}

// New returns a new API context.
func New(services *service.Service) *Context {
	return &Context{
		Services: services,
	}
}
