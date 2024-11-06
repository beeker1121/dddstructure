package transaction

import (
	"database/sql"
	"testing"

	"dddstructure/proto"
	"dddstructure/service"
	"dddstructure/storage/mock"
)

func TestProcess(t *testing.T) {
	// Create a new mock storage implementation.
	store := mock.New(&sql.DB{})

	// Create a new service.
	serv := service.New(store)

	// Create a user.
	u, err := serv.User.Create(&proto.UserCreateParams{
		Email:    "johndoe@fluidpay.com",
		Password: "TestPassword123",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Process a transaction.
	tx, err := serv.Transaction.Process(&proto.TransactionProcessParams{
		UserID: 1,
		Type:   "sale",
		Amount: 100,
		PaymentMethod: proto.TransactionPaymentMethod{
			Card: &proto.TransactionPaymentMethodCard{
				Number:         "4111111111111111",
				ExpirationDate: "1125",
			},
		},
	})
	if err != nil {
		panic(err)
	}

	// Check transaction.
	if tx.ID != 1 {
		t.Errorf("Expected transaction ID to be '%d', got '%d'", 1, tx.ID)
	}
	if tx.UserID != u.ID {
		t.Errorf("Expected transaction user ID to be '%d', got '%d'", u.ID, tx.UserID)
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
