package user

import (
	"fmt"
	"net/http"
	"strconv"

	apictx "dddstructure/cmd/api/context"
	"dddstructure/cmd/api/response"

	"github.com/beeker1121/httprouter"
)

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func New(ac *apictx.Context, router *httprouter.Router) {
	// Handle the routes.
	router.GET("/api/v1/user/:id", HandleGetUser(ac))
}

func HandleGetUser(ac *apictx.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the user ID.
		var id uint
		idu64, err := strconv.ParseUint(httprouter.GetParam(r, "id"), 10, 32)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		id = uint(idu64)

		// Get the user.
		servu, err := ac.Service.User.GetByID(id)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		// Map to API user response.
		u := &User{
			ID:       servu.ID,
			Username: servu.Username,
			Email:    servu.Email,
		}

		// Respond with JSON.
		if err := response.JSON(w, true, u); err != nil {
			fmt.Printf("error in handler: %v\n", err)
		}
	}
}
