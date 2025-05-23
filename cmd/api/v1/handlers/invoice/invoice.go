package invoice

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	router.GET("/api/v1/public/invoice/:hash", HandleGetPublicInvoice(ac))
	router.POST("/api/v1/public/invoice/:hash/pay", HandlePayInvoice(ac))
	router.GET("/api/v1/invoice/:id", auth.AuthenticateEndpoint(ac, HandleGetInvoice(ac)))
	router.POST("/api/v1/invoice/:id", auth.AuthenticateEndpoint(ac, HandlePostUpdate(ac)))
	router.DELETE("/api/v1/invoice/:id", auth.AuthenticateEndpoint(ac, HandleDelete(ac)))
}

// BillTo defines the billing information.
type BillTo struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Company      string `json:"company"`
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postal_code"`
	Country      string `json:"country"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
}

// PayTo defines the payee information.
type PayTo struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Company      string `json:"company"`
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postal_code"`
	Country      string `json:"country"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
}

// LineItems defines a line item.
type LineItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Quantity    uint   `json:"quantity"`
	Price       uint   `json:"price"`
}

// DueDate defines an invoice due date.
type DueDate struct {
	time.Time
}

// UnmarshalJSON implements the json.UnmarshalJSON interface.
func (dd *DueDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}

	dd.Time = t

	return nil
}

// MarshalJSON implements the json.MarshalJSON interface.
func (dd DueDate) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, dd.Format("2006-01-02"))), nil
}

// Invoice defines an invoice.
type Invoice struct {
	ID             uint                         `json:"id"`
	UserID         uint                         `json:"user_id"`
	PublicHash     string                       `json:"public_hash"`
	InvoiceNumber  string                       `json:"invoice_number"`
	PONumber       string                       `json:"po_number"`
	Currency       string                       `json:"currency"`
	DueDate        DueDate                      `json:"due_date"`
	Message        string                       `json:"message"`
	BillTo         BillTo                       `json:"bill_to"`
	PayTo          PayTo                        `json:"pay_to"`
	LineItems      []LineItem                   `json:"line_items"`
	PaymentMethods []proto.InvoicePaymentMethod `json:"payment_methods"`
	TaxRate        string                       `json:"tax_rate"`
	AmountDue      uint                         `json:"amount_due"`
	AmountPaid     uint                         `json:"amount_paid"`
	Status         string                       `json:"status"`
	CreatedAt      time.Time                    `json:"created_at"`
}

// RequestPost defines the request data for the HandlePost handler.
type RequestPost struct {
	InvoiceNumber  string                       `json:"invoice_number"`
	PONumber       string                       `json:"po_number"`
	Currency       string                       `json:"currency"`
	DueDate        DueDate                      `json:"due_date"`
	Message        string                       `json:"message"`
	BillTo         BillTo                       `json:"bill_to"`
	PayTo          PayTo                        `json:"pay_to"`
	LineItems      []LineItem                   `json:"line_items"`
	PaymentMethods []proto.InvoicePaymentMethod `json:"payment_methods"`
	TaxRate        string                       `json:"tax_rate"`
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

		// Handle line items.
		lineItems := []proto.InvoiceLineItem{}
		for _, v := range req.LineItems {
			lineItem := proto.InvoiceLineItem{
				Name:        v.Name,
				Description: v.Description,
				Quantity:    v.Quantity,
				Price:       v.Price,
			}

			lineItems = append(lineItems, lineItem)
		}

		// Create the invoice.
		invoice, err := ac.Service.Invoice.Create(&proto.InvoiceCreateParams{
			UserID:        user.ID,
			InvoiceNumber: req.InvoiceNumber,
			PONumber:      req.PONumber,
			Currency:      req.Currency,
			DueDate:       time.Time(req.DueDate.Time),
			Message:       req.Message,
			BillTo: proto.InvoiceBillTo{
				FirstName:    req.BillTo.FirstName,
				LastName:     req.BillTo.LastName,
				Company:      req.BillTo.Company,
				AddressLine1: req.BillTo.AddressLine1,
				AddressLine2: req.BillTo.AddressLine2,
				City:         req.BillTo.City,
				State:        req.BillTo.State,
				PostalCode:   req.BillTo.PostalCode,
				Country:      req.BillTo.Country,
				Email:        req.BillTo.Email,
				Phone:        req.BillTo.Phone,
			},
			PayTo: proto.InvoicePayTo{
				FirstName:    req.PayTo.FirstName,
				LastName:     req.PayTo.LastName,
				Company:      req.PayTo.Company,
				AddressLine1: req.PayTo.AddressLine1,
				AddressLine2: req.PayTo.AddressLine2,
				City:         req.PayTo.City,
				State:        req.PayTo.State,
				PostalCode:   req.PayTo.PostalCode,
				Country:      req.PayTo.Country,
				Email:        req.PayTo.Email,
				Phone:        req.PayTo.Phone,
			},
			LineItems:      lineItems,
			PaymentMethods: req.PaymentMethods,
			TaxRate:        req.TaxRate,
		})
		if pes, ok := err.(*serverrors.ParamErrors); ok && err != nil {
			errors.Params(ac.Logger, w, http.StatusBadRequest, pes)
			return
		} else if err != nil {
			ac.Logger.Error("invoice.Create() error",
				slog.Any("error", err))
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Create a new Result.
		result := ResultPost{
			Data: protoToInvoice(invoice),
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
	Data  []Invoice `json:"data"`
	Meta  Meta      `json:"meta"`
	Links Links     `json:"links"`
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

		// Handle created at start.
		if createdAtStartqs, ok := r.URL.Query()["created_at_start"]; ok && len(createdAtStartqs) == 1 {
			if params.CreatedAt == nil {
				params.CreatedAt = &proto.InvoiceGetParamsCreatedAt{}
			}

			t, err := time.Parse("2006-01-02 15:04:05", createdAtStartqs[0])
			if err != nil {
				errs.Add(errors.New(http.StatusBadRequest, "created_at_start", "invalid created at start date"))
			} else {
				params.CreatedAt.StartDate = &t
			}
		}

		// Handle created at end.
		if createdAtEndqs, ok := r.URL.Query()["created_at_end"]; ok && len(createdAtEndqs) == 1 {
			if params.CreatedAt == nil {
				params.CreatedAt = &proto.InvoiceGetParamsCreatedAt{}
			}

			t, err := time.Parse("2006-01-02 15:04:05", createdAtEndqs[0])
			if err != nil {
				errs.Add(errors.New(http.StatusBadRequest, "created_at_end", "invalid created at end date"))
			} else {
				params.CreatedAt.EndDate = &t
			}
		}

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
			ac.Logger.Error("invoice.Get() service error",
				slog.Any("error", err))
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Get invoices count.
		invoicesCount, err := ac.Service.Invoice.GetCount(params)
		if pes, ok := err.(*serverrors.ParamErrors); ok && err != nil {
			errors.Params(ac.Logger, w, http.StatusBadRequest, pes)
			return
		} else if err != nil {
			ac.Logger.Error("invoice.GetCount() service error",
				slog.Any("error", err))
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Create a new Result.
		result := ResultGet{
			Data: []Invoice{},
			Meta: Meta{
				Offset: params.Offset,
				Limit:  params.Limit,
				Total:  invoicesCount,
			},
			Links: Links{},
		}

		// Loop through the invoices.
		for _, i := range invoices {
			result.Data = append(result.Data, protoToInvoice(i))
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
			ac.Logger.Error("response.JSON() error",
				slog.Any("error", err))
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}
	}
}

// HandleGetPublicInvoice handles the /api/v1/invoice/public/:hash GET route of
// the API.
func HandleGetPublicInvoice(ac *apictx.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the hash.
		hash := httprouter.GetParam(r, "hash")

		// Get the invoice.
		invoice, err := ac.Service.Invoice.GetByPublicHash(hash)
		if err == serverrors.ErrInvoiceNotFound {
			errors.Default(ac.Logger, w, errors.New(http.StatusNotFound, "", err.Error()))
			return
		} else if err != nil {
			ac.Logger.Error("invoice.GetByPublicHash() service error",
				slog.Any("error", err))
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Create a new result.
		result := ResultGetInvoice{
			Data: protoToInvoice(invoice),
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

// RequestPayInvoice defines the request data for the HandlePayInvoice handler.
type RequestPayInvoice struct {
	Amount *uint `json:"amount"`
}

// ResultPayInvoice defines the response data for the HandlePayInvoice handler.
type ResultPayInvoice struct {
	Data Invoice `json:"data"`
}

// HandlePayInvoice handles the /api/v1/invoice/:id POST route of the API.
func HandlePayInvoice(ac *apictx.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the parameters from the request body.
		var req RequestPayInvoice
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			errors.Default(ac.Logger, w, errors.ErrBadRequest)
			return
		}

		// Try to get the invoice public hash.
		hash := httprouter.GetParam(r, "hash")

		// Get the invoice.
		invoice, err := ac.Service.Invoice.GetByPublicHash(hash)
		if err == serverrors.ErrInvoiceNotFound {
			errors.Default(ac.Logger, w, errors.New(http.StatusNotFound, "", err.Error()))
			return
		} else if err != nil {
			ac.Logger.Error("invoice.GetByPublicHash() service error",
				slog.Any("error", err))
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Pay the invoice.
		invoice, err = ac.Service.Invoice.Pay(invoice.ID, &proto.InvoicePayParams{
			Amount: *req.Amount,
		})
		if pes, ok := err.(*serverrors.ParamErrors); ok && err != nil {
			errors.Params(ac.Logger, w, http.StatusBadRequest, pes)
			return
		} else if err == serverrors.ErrInvoiceStatusNotPending {
			errors.Default(ac.Logger, w, errors.New(http.StatusBadRequest, "", err.Error()))
			return
		} else if err != nil {
			ac.Logger.Error("invoice.Pay() error",
				slog.Any("error", err))
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Create a new Result.
		result := ResultPayInvoice{
			Data: protoToInvoice(invoice),
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

// ResultGetInvoice defines the response data for the HandleGetInvoice handler.
type ResultGetInvoice struct {
	Data Invoice `json:"data"`
}

// HandleGetInvoice handles the /api/v1/invoice/:id GET route of the API.
func HandleGetInvoice(ac *apictx.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Try to get the invoice ID.
		var id uint
		id64, err := strconv.ParseInt(httprouter.GetParam(r, "id"), 10, 32)
		if err != nil {
			errors.Default(ac.Logger, w, errors.ErrBadRequest)
			return
		}
		id = uint(id64)

		// Get this user from the request context.
		user, err := auth.GetUserFromRequest(r)
		if err != nil {
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Get the invoice.
		invoice, err := ac.Service.Invoice.GetByIDAndUserID(id, user.ID)
		if err == serverrors.ErrInvoiceNotFound {
			errors.Default(ac.Logger, w, errors.New(http.StatusNotFound, "", err.Error()))
			return
		} else if err != nil {
			ac.Logger.Error("invoice.GetByIDAndUserID() service error",
				slog.Any("error", err))
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Create a new result.
		result := ResultGetInvoice{
			Data: protoToInvoice(invoice),
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

// BillToUpdate defines the billing information for update.
type BillToUpdate struct {
	FirstName    *string `json:"first_name"`
	LastName     *string `json:"last_name"`
	Company      *string `json:"company"`
	AddressLine1 *string `json:"address_line_1"`
	AddressLine2 *string `json:"address_line_2"`
	City         *string `json:"city"`
	State        *string `json:"state"`
	PostalCode   *string `json:"postal_code"`
	Country      *string `json:"country"`
	Email        *string `json:"email"`
	Phone        *string `json:"phone"`
}

// PayToUpdate defines the payee information for update.
type PayToUpdate struct {
	FirstName    *string `json:"first_name"`
	LastName     *string `json:"last_name"`
	Company      *string `json:"company"`
	AddressLine1 *string `json:"address_line_1"`
	AddressLine2 *string `json:"address_line_2"`
	City         *string `json:"city"`
	State        *string `json:"state"`
	PostalCode   *string `json:"postal_code"`
	Country      *string `json:"country"`
	Email        *string `json:"email"`
	Phone        *string `json:"phone"`
}

// RequestPostUpdate defines the request data for the HandlePostUpdate handler.
type RequestPostUpdate struct {
	InvoiceNumber  *string                       `json:"invoice_number"`
	PONumber       *string                       `json:"po_number"`
	Currency       *string                       `json:"currency"`
	DueDate        *DueDate                      `json:"due_date"`
	Message        *string                       `json:"message"`
	BillTo         *BillToUpdate                 `json:"bill_to"`
	PayTo          *PayToUpdate                  `json:"pay_to"`
	LineItems      *[]LineItem                   `json:"line_items"`
	PaymentMethods *[]proto.InvoicePaymentMethod `json:"payment_methods"`
	TaxRate        *string                       `json:"tax_rate"`
	AmountDue      *uint                         `json:"amount_due"`
}

// ResultPostUpdate defines the response data for the HandlePostUpdate handler.
type ResultPostUpdate struct {
	Data Invoice `json:"data"`
}

// HandlePostUpdate handles the /api/v1/invoice/:id POST route of the API.
func HandlePostUpdate(ac *apictx.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the parameters from the request body.
		var req RequestPostUpdate
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			errors.Default(ac.Logger, w, errors.ErrBadRequest)
			return
		}

		// Try to get the invoice ID.
		var id uint
		id64, err := strconv.ParseInt(httprouter.GetParam(r, "id"), 10, 32)
		if err != nil {
			errors.Default(ac.Logger, w, errors.ErrBadRequest)
			return
		}
		id = uint(id64)

		// Get this user from the request context.
		user, err := auth.GetUserFromRequest(r)
		if err != nil {
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Handle invoice update params.
		params := &proto.InvoiceUpdateParams{
			ID:             &id,
			UserID:         &user.ID,
			InvoiceNumber:  req.InvoiceNumber,
			PONumber:       req.PONumber,
			Currency:       req.Currency,
			Message:        req.Message,
			PaymentMethods: req.PaymentMethods,
			TaxRate:        req.TaxRate,
		}

		// Handle due date.
		if req.DueDate != nil {
			t := time.Time(req.DueDate.Time)
			params.DueDate = &t
		}

		if req.BillTo != nil {
			params.BillTo = &proto.InvoiceBillToUpdate{
				FirstName:    req.BillTo.FirstName,
				LastName:     req.BillTo.LastName,
				Company:      req.BillTo.Company,
				AddressLine1: req.BillTo.AddressLine1,
				AddressLine2: req.BillTo.AddressLine2,
				City:         req.BillTo.City,
				State:        req.BillTo.State,
				PostalCode:   req.BillTo.PostalCode,
				Country:      req.BillTo.Country,
				Email:        req.BillTo.Email,
				Phone:        req.BillTo.Phone,
			}
		}

		if req.PayTo != nil {
			params.PayTo = &proto.InvoicePayToUpdate{
				FirstName:    req.PayTo.FirstName,
				LastName:     req.PayTo.LastName,
				Company:      req.PayTo.Company,
				AddressLine1: req.PayTo.AddressLine1,
				AddressLine2: req.PayTo.AddressLine2,
				City:         req.PayTo.City,
				State:        req.PayTo.State,
				PostalCode:   req.PayTo.PostalCode,
				Country:      req.PayTo.Country,
				Email:        req.PayTo.Email,
				Phone:        req.PayTo.Phone,
			}
		}

		// Handle line items.
		if req.LineItems != nil {
			lineItems := []proto.InvoiceLineItem{}
			for _, v := range *req.LineItems {
				lineItem := proto.InvoiceLineItem{
					Name:        v.Name,
					Description: v.Description,
					Quantity:    v.Quantity,
					Price:       v.Price,
				}

				lineItems = append(lineItems, lineItem)
			}

			params.LineItems = &lineItems
		}

		// Update the invoice.
		invoice, err := ac.Service.Invoice.UpdateForUser(params)
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
		result := ResultPostUpdate{
			Data: protoToInvoice(invoice),
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

// HandleDelete handles the /api/v1/invoice/:id DELETE route of the API.
func HandleDelete(ac *apictx.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Try to get the invoice ID.
		var id uint
		id64, err := strconv.ParseInt(httprouter.GetParam(r, "id"), 10, 32)
		if err != nil {
			errors.Default(ac.Logger, w, errors.ErrBadRequest)
			return
		}
		id = uint(id64)

		// Get this user from the request context.
		user, err := auth.GetUserFromRequest(r)
		if err != nil {
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Get the invoice.
		_, err = ac.Service.Invoice.GetByIDAndUserID(id, user.ID)
		if err == serverrors.ErrInvoiceNotFound {
			errors.Default(ac.Logger, w, errors.New(http.StatusNotFound, "", err.Error()))
			return
		} else if err != nil {
			ac.Logger.Error("invoice.GetByIDAndUserID() service error",
				slog.Any("error", err))
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Delete the invoice.
		err = ac.Service.Invoice.Delete(id)
		if err == serverrors.ErrInvoiceNotFound {
			errors.Default(ac.Logger, w, errors.New(http.StatusNotFound, "", err.Error()))
			return
		} else if err != nil {
			ac.Logger.Error("invoice.GetByIDAndUserID() service error",
				slog.Any("error", err))
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// protoToInvoice handles mapping a proto invoice type to the response invoice
// type.
func protoToInvoice(i *proto.Invoice) Invoice {
	// Handle line items.
	lineItems := []LineItem{}
	for _, v := range i.LineItems {
		lineItem := LineItem{
			Name:        v.Name,
			Description: v.Description,
			Quantity:    v.Quantity,
			Price:       v.Price,
		}

		lineItems = append(lineItems, lineItem)
	}

	return Invoice{
		ID:            i.ID,
		UserID:        i.UserID,
		PublicHash:    i.PublicHash,
		InvoiceNumber: i.InvoiceNumber,
		PONumber:      i.PONumber,
		Currency:      i.Currency,
		DueDate:       DueDate{i.DueDate},
		Message:       i.Message,
		BillTo: BillTo{
			FirstName:    i.BillTo.FirstName,
			LastName:     i.BillTo.LastName,
			Company:      i.BillTo.Company,
			AddressLine1: i.BillTo.AddressLine1,
			AddressLine2: i.BillTo.AddressLine2,
			City:         i.BillTo.City,
			State:        i.BillTo.State,
			PostalCode:   i.BillTo.PostalCode,
			Country:      i.BillTo.Country,
			Email:        i.BillTo.Email,
			Phone:        i.BillTo.Phone,
		},
		PayTo: PayTo{
			FirstName:    i.PayTo.FirstName,
			LastName:     i.PayTo.LastName,
			Company:      i.PayTo.Company,
			AddressLine1: i.PayTo.AddressLine1,
			AddressLine2: i.PayTo.AddressLine2,
			City:         i.PayTo.City,
			State:        i.PayTo.State,
			PostalCode:   i.PayTo.PostalCode,
			Country:      i.PayTo.Country,
			Email:        i.PayTo.Email,
			Phone:        i.PayTo.Phone,
		},
		LineItems:      lineItems,
		PaymentMethods: i.PaymentMethods,
		TaxRate:        i.TaxRate,
		AmountDue:      i.AmountDue,
		AmountPaid:     i.AmountPaid,
		Status:         i.Status,
		CreatedAt:      i.CreatedAt,
	}
}
