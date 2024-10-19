package transaction

import (
	"database/sql"
	"testing"

	"dddstructure/proto"
	"dddstructure/service"
	"dddstructure/storage/mysql"
)

func TestProcess(t *testing.T) {
	// Create a new MySQL storage implementation.
	store := mysql.New(&sql.DB{})

	// Create a new service.
	serv := service.New(store)

	// Create a merchant.
	m, err := serv.Merchant.Create(&proto.Merchant{
		Name:  "John Doe",
		Email: "johndoe@fluidpay.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Process a transaction.
	tx, err := serv.Transaction.Process(&proto.Transaction{
		MerchantID:     m.ID,
		Type:           "sale",
		CardType:       "visa",
		AmountCaptured: 100,
	})
	if err != nil {
		panic(err)
	}

	// Check transaction.
	if tx.ID != 1 {
		t.Errorf("Expected transaction ID to be '%d', got '%d'", 1, tx.ID)
	}
	if tx.MerchantID != m.ID {
		t.Errorf("Expected transaction merchant ID to be '%d', got '%d'", m.ID, tx.MerchantID)
	}
	if tx.Type != "sale" {
		t.Errorf("Expected transaction type to be '%s', got '%s'", "sale", tx.Type)
	}
	if tx.CardType != "visa" {
		t.Errorf("Expected transaction card type to be '%s', got '%s'", "visa", tx.CardType)
	}
	if tx.AmountCaptured != 100 {
		t.Errorf("Expected transaction amount captured to be '%d', got '%d'", 100, tx.AmountCaptured)
	}
}
