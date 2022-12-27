package merchant

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"dddstructure/cmd/api/response"
	msctx "dddstructure/cmd/microservice/merchant/context"
	"dddstructure/proto"

	"github.com/beeker1121/httprouter"
)

func New(mc *msctx.Context, router *httprouter.Router) {
	// Handle the routes.
	router.GET("/rest/merchant/:id", HandleGetMerchant(mc))
	router.POST("/rest/merchant", HandlePostMerchant(mc))
}

func HandleGetMerchant(mc *msctx.Context) http.HandlerFunc {
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
		servm, err := mc.Service.Merchant.GetByID(id)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		// Map to API user response.
		m := &proto.Merchant{
			ID:    servm.ID,
			Name:  servm.Name,
			Email: servm.Email,
		}

		// Respond with JSON.
		if err := response.JSON(w, true, m); err != nil {
			fmt.Printf("error in handler: %v\n", err)
		}
	}
}

func HandlePostMerchant(mc *msctx.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request body to Merchant type.
		var m proto.Merchant
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		// Create a new merchant.
		servm, err := mc.Service.Merchant.Create(&proto.Merchant{
			ID:    m.ID,
			Name:  m.Name,
			Email: m.Email,
		})
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		// Map to API merchant response.
		resm := &proto.Merchant{
			ID:    servm.ID,
			Name:  servm.Name,
			Email: servm.Email,
		}

		// Respond with JSON.
		if err := response.JSON(w, true, resm); err != nil {
			fmt.Printf("error in handler: %v\n", err)
		}
	}
}
