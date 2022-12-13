package merchant

import (
	"fmt"
	"net/http"
	"strconv"

	apictx "dddstructure/cmd/api/context"
	"dddstructure/cmd/api/response"

	"github.com/beeker1121/httprouter"
)

type Merchant struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func New(ac *apictx.Context, router *httprouter.Router) {
	// Handle the routes.
	router.GET("/api/v1/merchant/:id", HandleGetMerchant(ac))
}

func HandleGetMerchant(ac *apictx.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the merchant ID.
		var id uint
		idu64, err := strconv.ParseUint(httprouter.GetParam(r, "id"), 10, 32)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		id = uint(idu64)

		// Get the merchant.
		servm, err := ac.Service.Merchant.GetByID(id)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		// Map to API merchant response.
		m := &Merchant{
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
