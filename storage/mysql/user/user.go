package user

import (
	"database/sql"
	"errors"
	"fmt"

	"dddstructure/storage/user"
)

// userMap acts as a mock MySQL database for users.
var userMap map[uint]*user.User = make(map[uint]*user.User)

// Database defines the database.
type Database struct {
	db *sql.DB
}

// New creates a new database.
func New(db *sql.DB) *Database {
	return &Database{
		db: db,
	}
}

// Create creates a new user.
func (db *Database) Create(params *user.CreateParams) (*user.User, error) {
	u := &user.User{
		ID:    params.ID,
		Name:  params.Name,
		Email: params.Email,
	}

	userMap[u.ID] = u

	fmt.Println("Created user and added to MySQL database...")

	return u, nil
}

// GetByID gets a user by the given ID.
func (db *Database) GetByID(id uint) (*user.User, error) {
	u, ok := userMap[id]
	if !ok {
		return nil, errors.New("could not find merchant")
	}

	fmt.Println("Got merchant from MySQL database...")

	return u, nil
}
