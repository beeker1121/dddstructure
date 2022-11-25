package v1

import (
	apictx "dddstructure/cmd/api/context"
	"dddstructure/cmd/api/v1/handlers/merchant"

	"github.com/beeker1121/httprouter"
)

// New creates a new API v1 application. All of the necessary routes for
// v1 of the API will be created on the given router, which should then be
// used to create the web server.
func New(ac *apictx.Context, router *httprouter.Router) {
	// Create all of the API v1 routes.
	merchant.New(ac, router)
}
