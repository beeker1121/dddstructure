package user

import (
	"encoding/json"
	"log/slog"
	"net/http"

	apictx "dddstructure/cmd/api/context"
	"dddstructure/cmd/api/errors"
	"dddstructure/cmd/api/middleware/auth"
	"dddstructure/cmd/api/response"
	"dddstructure/proto"
	serverrors "dddstructure/service/errors"

	"github.com/beeker1121/httprouter"
)

// New creates the routes for the user endpoints of the API.
func New(ac *apictx.Context, router *httprouter.Router) {
	// Handle the routes.
	router.GET("/api/v1/user", auth.AuthenticateEndpoint(ac, HandleGet(ac)))
	router.POST("/api/v1/user", auth.AuthenticateEndpoint(ac, HandlePost(ac)))
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
			ac.Logger.Error("user.GetByID() service error",
				slog.Any("error", err))
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
			ac.Logger.Error("response.JSON() error",
				slog.Any("error", err))
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}
	}
}

// RequestPost defines the request data for the HandlePost handler.
type RequestPost struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

// ResultPost defines the response data for the HandlePost handler.
type ResultPost struct {
	Data User `json:"data"`
}

// HandlePost handles the /api/v1/user POST route of the API.
func HandlePost(ac *apictx.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the parameters from the request body.
		var req RequestPost
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			errors.Default(ac.Logger, w, errors.ErrBadRequest)
			return
		}

		// Get this user from the request context.
		user, err := auth.GetUserFromRequest(r)
		if err != nil {
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Update the user.
		user, err = ac.Service.User.Update(&proto.UserUpdateParams{
			ID:       &user.ID,
			Email:    req.Email,
			Password: req.Password,
		})
		if pes, ok := err.(*serverrors.ParamErrors); ok && err != nil {
			errors.Params(ac.Logger, w, http.StatusBadRequest, pes)
			return
		} else if err != nil {
			ac.Logger.Error("user.Update() error",
				slog.Any("error", err))
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Create a new Result.
		result := ResultPost{
			Data: User{
				ID:    user.ID,
				Email: user.Email,
			},
		}

		// Respond with JSON.
		if err := response.JSON(w, true, result); err != nil {
			ac.Logger.Error("response.JSON() error",
				slog.Any("error", err))
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}
	}
}
