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

	// Create a user.
	u, err := serv.User.Create(&proto.User{
		Username: "johndoe",
		Email:    "johndoe@fluidpay.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check user.
	if u.ID != 1 {
		t.Errorf("Expected user ID to be '%d', got '%d'", 2, u.ID)
	}
	if u.Email != "johndoe@fluidpay.com" {
		t.Errorf("Expected user email to be '%s', got '%s'", "johndoe@fluidpay.com", u.Email)
	}
	if u.Username != "johndoe" {
		t.Errorf("Expected user username to be '%s', got '%s'", "johndoe", u.Username)
	}
}
