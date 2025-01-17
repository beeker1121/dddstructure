package login

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

// New creates the routes for the login endpoints of the API.
func New(ac *apictx.Context, router *httprouter.Router) {
	// Handle the routes.
	router.POST("/api/v1/login", HandlePost(ac))
}

// RequestPost defines the request data for the HandlePost handler.
type RequestPost struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ResultPost defines the response data for the HandlePost handler.
type ResultPost struct {
	Data string `json:"data"`
}

// HandlePost handles the /api/v1/login POST route of the API.
func HandlePost(ac *apictx.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the parameters from the request body.
		var req RequestPost
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			errors.Default(ac.Logger, w, errors.ErrBadRequest)
			return
		}

		// Try to log in the user.
		user, err := ac.Service.User.Login(&proto.UserLoginParams{
			Email:    req.Email,
			Password: req.Password,
		})
		if pes, ok := err.(*serverrors.ParamErrors); ok && err != nil {
			errors.Params(ac.Logger, w, http.StatusBadRequest, pes)
			return
		} else if err == serverrors.ErrUserInvalidLogin {
			errors.Default(ac.Logger, w, errors.New(http.StatusUnauthorized, "", err.Error()))
			return
		} else if err != nil {
			ac.Logger.Error("user.Login() service error",
				slog.Any("error", err))
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Issue a new JWT for this user.
		token, err := auth.NewJWT(ac, user.Password, user.ID)
		if err != nil {
			ac.Logger.Error("auth.NewJWT() error",
				slog.Any("error", err))
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Create a new Result.
		result := ResultPost{
			Data: token,
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
