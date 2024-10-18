package merchant

import (
	"database/sql"
	"testing"

	"dddstructure/proto"
	"dddstructure/service"
	"dddstructure/storage/mysql"
)

func TestCreate(t *testing.T) {
	// Create a new MySQL storage implementation.
	store := mysql.New(&sql.DB{})

	// Create a new service.
	serv := service.New(store)

	// Register dependencies.
	dep.RegisterMerchant(serv.Merchant)
	dep.RegisterUser(serv.User)
	dep.RegisterInvoice(serv.Invoice)
	dep.RegisterTransaction(serv.Transaction)

	// Create a merchant.
	m, err := serv.Merchant.Create(&proto.Merchant{
		Name:  "John Doe",
		Email: "johndoe@fluidpay.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check merchant.
	if m.ID != 1 {
		t.Errorf("Expected merchant ID to be '%d', got '%d'", 1, m.ID)
	}
	if m.Name != "John Doe" {
		t.Errorf("Expected merchant name to be '%s', got '%s'", "John Doe", m.Name)
	}
	if m.Email != "johndoe@fluidpay.com" {
		t.Errorf("Expected merchant email to be '%s', got '%s'", "johndoe@fluidpay.com", m.Email)
	}
}
