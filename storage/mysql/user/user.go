package user

import (
	"database/sql"
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
	m := &user.User{
		ID:    params.ID,
		Name:  params.Name,
		Email: params.Email,
	}

	userMap[m.ID] = m

	fmt.Println("Created user and added to MySQL database...")

	return m, nil
}

// GetByID gets a user by the given ID.
func (db *Database) GetByID(id uint) (*user.User, error) {
	u := &user.User{
		ID:    1,
		Name:  "John Doe",
		Email: "johndoe@gmail.com",
	}

	return u, nil
}
