package transaction

import (
	"encoding/json"
	"net/http"

	apictx "dddstructure/cmd/api/context"
	"dddstructure/cmd/api/errors"
	"dddstructure/cmd/api/middleware/auth"
	"dddstructure/cmd/api/response"
	"dddstructure/proto"
	serverrors "dddstructure/service/errors"

	"github.com/beeker1121/httprouter"
)

// New creates the routes for the transaction endpoints of the API.
func New(ac *apictx.Context, router *httprouter.Router) {
	// Handle the routes.
	router.POST("/api/v1/transaction", auth.AuthenticateEndpoint(ac, HandlePost(ac)))
}

// Transaction defines a transaction.
type Transaction struct {
	ID             uint   `json:"id"`
	UserID         uint   `json:"user_id"`
	Type           string `json:"type"`
	CardType       string `json:"card_type"`
	AmountCaptured uint   `json:"amount_captured"`
	InvoiceID      uint   `json:"invoice_id"`
	Status         string `json:"status"`
}

// PaymentMethod defines the transaction payment method.
type PaymentMethod struct {
	Card *PaymentMethodCard `json:"card"`
}

// PaymentMethodCard defines the card payment method.
type PaymentMethodCard struct {
	Number         string `json:"number"`
	ExpirationDate string `json:"expiration_date"`
}

// RequestPost defines the request data for the HandlePost handler.
type RequestPost struct {
	Type          string        `json:"type"`
	Amount        uint          `json:"amount"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	InvoiceID     uint          `json:"invoice_id"`
}

// ResultPost defines the response data for the HandlePost handler.
type ResultPost struct {
	Data *Transaction `json:"data"`
}

// HandlePost handles the /api/v1/transaction POST route of the API.
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

		// Create the transaction payment method.
		paymentMethod := proto.TransactionPaymentMethod{}
		if req.PaymentMethod.Card != nil {
			paymentMethod.Card = &proto.TransactionPaymentMethodCard{
				Number:         req.PaymentMethod.Card.Number,
				ExpirationDate: req.PaymentMethod.Card.ExpirationDate,
			}
		}

		// Process the transaction.
		transaction, err := ac.Service.Transaction.Process(&proto.TransactionProcessParams{
			UserID:        user.ID,
			Type:          req.Type,
			Amount:        req.Amount,
			PaymentMethod: paymentMethod,
			InvoiceID:     req.InvoiceID,
		})
		if pes, ok := err.(*serverrors.ParamErrors); ok && err != nil {
			errors.Params(ac.Logger, w, http.StatusBadRequest, pes)
			return
		} else if err != nil {
			ac.Logger.Printf("transaction.Process() error: %s\n", err)
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Create a new Result.
		result := ResultPost{
			Data: &Transaction{
				ID:             transaction.ID,
				UserID:         transaction.UserID,
				Type:           transaction.Type,
				CardType:       transaction.CardType,
				AmountCaptured: transaction.AmountCaptured,
				InvoiceID:      transaction.InvoiceID,
				Status:         transaction.Status,
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
