package invoice

import (
	"encoding/json"
	"net/http"
	"strconv"

	apictx "dddstructure/cmd/api/context"
	"dddstructure/cmd/api/errors"
	"dddstructure/cmd/api/middleware/auth"
	"dddstructure/cmd/api/response"
	"dddstructure/proto"
	serverrors "dddstructure/service/errors"

	"github.com/beeker1121/httprouter"
)

// New creates the routes for the invoice endpoints of the API.
func New(ac *apictx.Context, router *httprouter.Router) {
	// Handle the routes.
	router.POST("/api/v1/invoice", auth.AuthenticateEndpoint(ac, HandlePost(ac)))
	router.GET("/api/v1/invoice", auth.AuthenticateEndpoint(ac, HandleGet(ac)))
}

// BillTo defines the billing information.
type BillTo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// PayTo defines the payee information.
type PayTo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Invoice defines an invoice.
type Invoice struct {
	ID         uint   `json:"id"`
	UserID     uint   `json:"user_id"`
	BillTo     BillTo `json:"bill_to"`
	PayTo      PayTo  `json:"pay_to"`
	AmountDue  uint   `json:"amount_due"`
	AmountPaid uint   `json:"amount_paid"`
	Status     string `json:"status"`
}

// RequestPost defines the request data for the HandlePost handler.
type RequestPost struct {
	BillTo    BillTo `json:"bill_to"`
	PayTo     PayTo  `json:"pay_to"`
	AmountDue uint   `json:"amount_due"`
}

// ResultPost defines the response data for the HandlePost handler.
type ResultPost struct {
	Data Invoice `json:"data"`
}

// HandlePost handles the /api/v1/invoice POST route of the API.
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

		// Create the invoice.
		invoice, err := ac.Service.Invoice.Create(&proto.InvoiceCreateParams{
			UserID: user.ID,
			BillTo: proto.InvoiceBillTo{
				FirstName: req.BillTo.FirstName,
				LastName:  req.BillTo.LastName,
			},
			PayTo: proto.InvoicePayTo{
				FirstName: req.PayTo.FirstName,
				LastName:  req.PayTo.LastName,
			},
			AmountDue: req.AmountDue,
		})
		if pes, ok := err.(*serverrors.ParamErrors); ok && err != nil {
			errors.Params(ac.Logger, w, http.StatusBadRequest, pes)
			return
		} else if err != nil {
			ac.Logger.Printf("invoice.Create() error: %s\n", err)
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Create a new Result.
		result := ResultPost{
			Data: Invoice{
				ID:     invoice.ID,
				UserID: invoice.UserID,
				BillTo: BillTo{
					FirstName: invoice.BillTo.FirstName,
					LastName:  invoice.BillTo.LastName,
				},
				PayTo: PayTo{
					FirstName: invoice.PayTo.FirstName,
					LastName:  invoice.PayTo.LastName,
				},
				AmountDue:  invoice.AmountDue,
				AmountPaid: invoice.AmountPaid,
				Status:     invoice.Status,
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

// Meta defines the response top level meta object.
type Meta struct {
	Offset uint `json:"offset"`
	Limit  uint `json:"limit"`
	Total  uint `json:"total"`
}

// Links defines the response top level links object.
type Links struct {
	Prev *string `json:"prev"`
	Next *string `json:"next"`
}

// ResultGet defines the response data for the HandleGet handler.
type ResultGet struct {
	Data  []*Invoice `json:"data"`
	Meta  Meta       `json:"meta"`
	Links Links      `json:"links"`
}

// HandleGet handles the /api/v1/invoice GET route of the API.
func HandleGet(ac *apictx.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get this user from the request context.
		user, err := auth.GetUserFromRequest(r)
		if err != nil {
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Create a new GetParams.
		params := &proto.InvoiceGetParams{
			UserID: &user.ID,
		}

		// Create a new API Errors.
		errs := &errors.Errors{}

		// Handle offset.
		if offsetqs, ok := r.URL.Query()["offset"]; ok && len(offsetqs) == 1 {
			offset64, err := strconv.ParseInt(offsetqs[0], 10, 32)
			if err != nil {
				errs.Add(errors.ErrOffsetInvalid)
			} else {
				params.Offset = uint(offset64)
			}
		} else {
			params.Offset = 0
		}

		// Handle limit.
		if limitqs, ok := r.URL.Query()["limit"]; ok && len(limitqs) == 1 {
			limit64, err := strconv.ParseInt(limitqs[0], 10, 32)
			if err != nil {
				errs.Add(errors.ErrLimitInvalid)
			} else {
				if uint(limit64) > ac.Config.LimitMax {
					errs.Add(errors.ErrLimitMax(uint(limit64), ac.Config.LimitMax))
				} else {
					params.Limit = uint(limit64)
				}
			}
		} else {
			params.Limit = ac.Config.LimitDefault
		}

		// Return if there were errors.
		if errs.Length() > 0 {
			errors.Multiple(ac.Logger, w, http.StatusBadRequest, errs)
			return
		}

		// Get invoices.
		invoices, err := ac.Service.Invoice.Get(params)
		if pes, ok := err.(*serverrors.ParamErrors); ok && err != nil {
			errors.Params(ac.Logger, w, http.StatusBadRequest, pes)
			return
		} else if err != nil {
			ac.Logger.Printf("invoice.Get() service error: %s\n", err)
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Get invoices count.
		invoicesCount, err := ac.Service.Invoice.GetCount(params)
		if pes, ok := err.(*serverrors.ParamErrors); ok && err != nil {
			errors.Params(ac.Logger, w, http.StatusBadRequest, pes)
			return
		} else if err != nil {
			ac.Logger.Printf("invoice.GetCount() service error: %s\n", err)
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Create a new Result.
		result := ResultGet{
			Data: []*Invoice{},
			Meta: Meta{
				Offset: params.Offset,
				Limit:  params.Limit,
				Total:  invoicesCount,
			},
			Links: Links{},
		}

		// Loop through the invoices.
		for _, i := range invoices {
			// Create a new invoice.
			invoice := &Invoice{
				ID:     i.ID,
				UserID: i.UserID,
				BillTo: BillTo{
					FirstName: i.BillTo.FirstName,
					LastName:  i.BillTo.LastName,
				},
				PayTo: PayTo{
					FirstName: i.PayTo.FirstName,
					LastName:  i.PayTo.LastName,
				},
				AmountDue:  i.AmountDue,
				AmountPaid: i.AmountPaid,
				Status:     i.Status,
			}

			result.Data = append(result.Data, invoice)
		}

		// Handle previous link.
		if params.Offset > 0 {
			limitstr := "&limit=" + strconv.FormatInt(int64(params.Limit), 10)

			offsetstr := "?offset="
			if params.Offset-params.Limit < 0 {
				offsetstr += "0"
			} else {
				offsetstr += strconv.FormatInt(int64(params.Offset-params.Limit), 10)
			}

			prev := "https://" + ac.Config.APIHost + "/api/v1/invoice" + offsetstr + limitstr
			result.Links.Prev = &prev
		}

		// Handle next link.
		if params.Offset+params.Limit < result.Meta.Total {
			offsetstr := "?offset=" + strconv.FormatInt(int64(params.Offset+params.Limit), 10)
			limitstr := "&limit=" + strconv.FormatInt(int64(params.Limit), 10)

			next := "https://" + ac.Config.APIHost + "/api/v1/invoice" + offsetstr + limitstr
			result.Links.Next = &next
		}

		// Respond with JSON.
		if err := response.JSON(w, true, result); err != nil {
			ac.Logger.Printf("response.JSON() error: %s\n", err)
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}
	}
}
