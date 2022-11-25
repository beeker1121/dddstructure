package merchant

import (
	"encoding/json"
	"net/http"
	"strconv"

	apictx "dddstructure/cmd/api/context"
	"dddstructure/service/merchant"

	"github.com/beeker1121/httprouter"
)

// New creates the routes for the todo endpoints of the API.
func New(ac *apictx.Context, router *httprouter.Router) {
	router.GET("/api/v1/merchant/:id", HandleGetMerchant(ac))
	router.POST("/api/v1/merchant", HandlePost(ac))
}

type GetMerchantResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// HandleGetMerchant handles the /api/v1/merchant/:id GET route of the API.
func HandleGetMerchant(ac *apictx.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the merchant ID.
		var id uint
		id64, err := strconv.ParseUint(httprouter.GetParam(r, "id"), 10, 32)
		if err != nil {
			w.Write([]byte("error"))
			return
		}
		id = uint(id64)

		// Get the merchant via services.
		m, err := ac.Services.Merchant.GetByID(id)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		// Map to response.
		res := GetMerchantResponse{
			ID:    m.ID,
			Name:  m.Name,
			Email: m.Email,
		}

		// Marshal response to JSON.
		resJSON, err := json.Marshal(res)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(resJSON)
	}
}

type PostRequest struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type PostResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// HandlePost handles the /api/v1/merchant POST route of the API.
func HandlePost(ac *apictx.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request.
		var postReq PostRequest
		if err := json.NewDecoder(r.Body).Decode(&postReq); err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		// Map request to the merchant create params.
		params := &merchant.CreateParams{
			ID:    postReq.ID,
			Name:  postReq.Name,
			Email: postReq.Email,
		}

		// Create a merchant.
		m, err := ac.Services.Merchant.Create(params)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		// Map to response.
		res := PostResponse{
			ID:    m.ID,
			Name:  m.Name,
			Email: m.Email,
		}

		// Marshal response to JSON.
		resJSON, err := json.Marshal(res)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(resJSON)
	}
}
