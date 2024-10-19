package user

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

	// Create a merchant.
	m, err := serv.Merchant.Create(&proto.Merchant{
		Name:  "John Doe",
		Email: "johndoe@fluidpay.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a user.
	u, err := serv.User.Create(&proto.User{
		AccountTypeID: m.ID,
		Username:      "johndoe",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check user.
	if u.ID != 2 {
		t.Errorf("Expected user ID to be '%d', got '%d'", 2, u.ID)
	}
	if u.AccountTypeID != m.ID {
		t.Errorf("Expected user account type ID to be '%d', got '%d'", m.ID, u.AccountTypeID)
	}
	if u.Username != "johndoe" {
		t.Errorf("Expected user username to be '%s', got '%s'", "johndoe", u.Username)
	}
}
