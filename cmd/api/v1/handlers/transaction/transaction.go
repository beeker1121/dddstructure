package transaction

import (
	"encoding/json"
	"fmt"
	"net/http"

	apictx "dddstructure/cmd/api/context"
	"dddstructure/cmd/api/response"
	"dddstructure/proto"

	"github.com/beeker1121/httprouter"
)

type Transaction struct {
	ID             uint
	MerchantID     uint
	Type           string
	ProcessorType  string
	CardType       string
	AmountCaptured uint
	InvoiceID      uint
}

func New(ac *apictx.Context, router *httprouter.Router) {
	// Handle the routes.
	router.GET("/api/v1/transaction", HandlePostTransaction(ac))
}

func HandlePostTransaction(ac *apictx.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request body to Transaction type.
		var t Transaction
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		// Process the transaction.
		servt, err := ac.Service.Transaction.Process(&proto.Transaction{
			ID:             t.ID,
			MerchantID:     t.MerchantID,
			Type:           t.Type,
			ProcessorType:  t.ProcessorType,
			CardType:       t.CardType,
			AmountCaptured: t.AmountCaptured,
			InvoiceID:      t.InvoiceID,
		})
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		// Map to API transaction response.
		rest := &Transaction{
			ID:             servt.ID,
			MerchantID:     servt.MerchantID,
			Type:           servt.Type,
			ProcessorType:  servt.ProcessorType,
			CardType:       servt.CardType,
			AmountCaptured: servt.AmountCaptured,
			InvoiceID:      servt.InvoiceID,
		}

		// Respond with JSON.
		if err := response.JSON(w, true, rest); err != nil {
			fmt.Printf("error in handler: %v\n", err)
		}
	}
}
