package invoice

import (
	"database/sql"
	"testing"

	"dddstructure/proto"
	"dddstructure/service"
	"dddstructure/storage/mysql"
)

func TestPay(t *testing.T) {
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

	// Create an invoice.
	i, err := serv.Invoice.Create(&proto.Invoice{
		UserID:     u.ID,
		BillTo:     "Joe Smith",
		PayTo:      "John Doe",
		AmountDue:  100,
		AmountPaid: 0,
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
	i, err = serv.Invoice.Pay(i.ID)
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
