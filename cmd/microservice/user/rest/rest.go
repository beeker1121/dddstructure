package rest

import (
	msctx "dddstructure/cmd/microservice/user/context"
	"dddstructure/cmd/microservice/user/rest/handlers/user"

	"github.com/beeker1121/httprouter"
)

func New(mc *msctx.Context, r *httprouter.Router) {
	user.New(mc, r)
}
