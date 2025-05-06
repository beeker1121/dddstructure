package user

import (
	"database/sql"
	"log/slog"
	"testing"

	"dddstructure/proto"
	"dddstructure/service"
	"dddstructure/storage/mock"
)

func TestCreate(t *testing.T) {
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

	// Check user.
	if u.ID != 1 {
		t.Errorf("Expected user ID to be '%d', got '%d'", 2, u.ID)
	}
	if u.Email != "johndoe@test.com" {
		t.Errorf("Expected user email to be '%s', got '%s'", "johndoe@test.com", u.Email)
	}
}
