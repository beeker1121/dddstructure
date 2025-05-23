package invoice

import (
	"database/sql"
	"log/slog"
	"testing"

	"dddstructure/proto"
	"dddstructure/service"
	"dddstructure/storage/mock"
)

func TestPay(t *testing.T) {
	// Create a new mock storage implementation.
	store := mock.New(&sql.DB{})

	// Create a new service.
	serv := service.New(store, &slog.Logger{})

	// Create a user.
	u, err := serv.User.Create(&proto.UserCreateParams{
		Email:    "johndoe@test.com",
		Password: "TestPassword123",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create an invoice.
	i, err := serv.Invoice.Create(&proto.InvoiceCreateParams{
		UserID: u.ID,
		BillTo: proto.InvoiceBillTo{
			FirstName: "John",
			LastName:  "Smith",
		},
		PayTo: proto.InvoicePayTo{
			FirstName: "John",
			LastName:  "Doe",
		},
		LineItems: []proto.InvoiceLineItem{
			{
				Quantity: 1,
				Price:    100,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check invoice.
	if i.AmountDue != 100 {
		t.Errorf("Expected invoice amount due to be '%d', got '%d'", 100, i.AmountDue)
	}
	if i.AmountPaid != 0 {
		t.Errorf("Expected invoice amount paid to be '%d', got '%d'", 0, i.AmountPaid)
	}
	if i.Status != "pending" {
		t.Errorf("Expected invoice status to be '%s', got '%s'", "pending", i.Status)
	}

	// Pay invoice.
	i, err = serv.Invoice.Pay(i.ID, &proto.InvoicePayParams{
		Amount: 100,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check invoice.
	if i.AmountDue != 0 {
		t.Errorf("Expected amount due to be '%d', got '%d'", 0, i.AmountDue)
	}
	if i.AmountPaid != 100 {
		t.Errorf("Expected amount paid to be '%d', got '%d'", 100, i.AmountPaid)
	}
	if i.Status != "paid" {
		t.Errorf("Expected status to be '%s', got '%s'", "paid", i.Status)
	}
}
