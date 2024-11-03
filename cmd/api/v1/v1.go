package v1

import (
	apictx "dddstructure/cmd/api/context"
	"dddstructure/cmd/api/v1/handlers/invoice"
	"dddstructure/cmd/api/v1/handlers/signup"
	"dddstructure/cmd/api/v1/handlers/user"

	"github.com/beeker1121/httprouter"
)

func New(ac *apictx.Context, r *httprouter.Router) {
	invoice.New(ac, r)
	signup.New(ac, r)
	user.New(ac, r)
}
