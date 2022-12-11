package v1

import (
	apictx "dddstructure/cmd/api/context"
	"dddstructure/cmd/api/v1/handlers/merchant"
	"dddstructure/cmd/api/v1/handlers/user"

	"github.com/beeker1121/httprouter"
)

func New(ac *apictx.Context, r *httprouter.Router) {
	merchant.New(ac, r)
	user.New(ac, r)
}
