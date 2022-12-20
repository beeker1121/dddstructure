package user

import (
	"fmt"
	"net/http"
	"strconv"

	"dddstructure/cmd/api/response"
	msctx "dddstructure/cmd/microservice/user/context"

	"github.com/beeker1121/httprouter"
)

type User struct {
	ID            uint   `json:"id"`
	AccountTypeID uint   `json:"account_type_id"`
	Username      string `json:"username"`
}

func New(mc *msctx.Context, router *httprouter.Router) {
	// Handle the routes.
	router.GET("/api/v1/user/:id", HandleGetUser(mc))
}

func HandleGetUser(mc *msctx.Context) http.HandlerFunc {
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
		servu, err := mc.Service.User.GetByID(id)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		// Map to API user response.
		u := &User{
			ID:            servu.ID,
			AccountTypeID: servu.AccountTypeID,
			Username:      servu.Username,
		}

		// Respond with JSON.
		if err := response.JSON(w, true, u); err != nil {
			fmt.Printf("error in handler: %v\n", err)
		}
	}
}
