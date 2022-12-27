package rest

import (
	msctx "dddstructure/cmd/microservice/merchant/context"
	"dddstructure/cmd/microservice/merchant/rest/handlers/merchant"

	"github.com/beeker1121/httprouter"
)

func New(mc *msctx.Context, r *httprouter.Router) {
	merchant.New(mc, r)
}
