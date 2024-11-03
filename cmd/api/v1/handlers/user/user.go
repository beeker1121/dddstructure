package user

import (
	"net/http"

	apictx "dddstructure/cmd/api/context"
	"dddstructure/cmd/api/errors"
	"dddstructure/cmd/api/middleware/auth"
	"dddstructure/cmd/api/response"

	"github.com/beeker1121/httprouter"
)

// New creates the routes for the user endpoints of the API.
func New(ac *apictx.Context, router *httprouter.Router) {
	// Handle the routes.
	router.GET("/api/v1/user", auth.AuthenticateEndpoint(ac, HandleGet(ac)))
}

// User defines a user.
type User struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

// ResultGet defines the response data for the HandleGet handler.
type ResultGet struct {
	Data User `json:"data"`
}

// HandleGet handles the /api/v1/user GET route of the API.
func HandleGet(ac *apictx.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get this user from the request context.
		user, err := auth.GetUserFromRequest(r)
		if err != nil {
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Get the user.
		serviceu, err := ac.Service.User.GetByID(user.ID)
		if err != nil {
			ac.Logger.Printf("user.GetByID() service error: %s\n", err)
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Create a new Result.
		result := ResultGet{
			Data: User{
				ID:    serviceu.ID,
				Email: serviceu.Email,
			},
		}

		// Respond with JSON.
		if err := response.JSON(w, true, result); err != nil {
			ac.Logger.Printf("response.JSON() error: %s\n", err)
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}
	}
}
