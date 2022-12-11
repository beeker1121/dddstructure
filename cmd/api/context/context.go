package context

import "dddstructure/service"

type Context struct {
	Service *service.Service
}

func New(s *service.Service) *Context {
	return &Context{
		Service: s,
	}
}
