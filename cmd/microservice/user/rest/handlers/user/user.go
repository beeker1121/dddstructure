package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"dddstructure/cmd/api/response"
	msctx "dddstructure/cmd/microservice/user/context"
	"dddstructure/proto"

	"github.com/beeker1121/httprouter"
)

func New(mc *msctx.Context, router *httprouter.Router) {
	// Handle the routes.
	router.GET("/rest/user/:id", HandleGetUser(mc))
	router.POST("/rest/user", HandlePostUser(mc))
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
		u := &proto.User{
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

func HandlePostUser(mc *msctx.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request body to User type.
		var u proto.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		// Create a new user.
		servu, err := mc.Service.User.Create(&proto.User{
			AccountTypeID: u.AccountTypeID,
			Username:      u.Username,
		})
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		// Map to API user response.
		resu := &proto.User{
			ID:            servu.ID,
			AccountTypeID: servu.AccountTypeID,
			Username:      servu.Username,
		}

		// Respond with JSON.
		if err := response.JSON(w, true, resu); err != nil {
			fmt.Printf("error in handler: %v\n", err)
		}
	}
}
