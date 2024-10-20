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

	// Create a user.
	u, err := serv.User.Create(&proto.User{
		Username: "johndoe",
		Email:    "johndoe@fluidpay.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Process a transaction.
	tx, err := serv.Transaction.Process(&proto.Transaction{
		UserID:         u.ID,
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
	if tx.UserID != u.ID {
		t.Errorf("Expected transaction merchant ID to be '%d', got '%d'", u.ID, tx.UserID)
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
