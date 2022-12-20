package context

import "dddstructure/cmd/microservice/user/service"

type Context struct {
	Service *service.Service
}

func New(s *service.Service) *Context {
	return &Context{
		Service: s,
	}
}
