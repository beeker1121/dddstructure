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
func (db *Database) Create(u *user.User) (*user.User, error) {
	use := &user.User{
		ID:            u.ID,
		AccountTypeID: u.AccountTypeID,
		Username:      u.Username,
	}

	userMap[use.ID] = use

	fmt.Println("Created user and added to MySQL database...")

	return use, nil
}

// GetByID gets an user by the given ID.
func (db *Database) GetByID(id uint) (*user.User, error) {
	m, ok := userMap[id]
	if !ok {
		return nil, errors.New("could not find user")
	}

	fmt.Println("Got user from MySQL database...")

	return m, nil
}
