package user

import (
	"database/sql"

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
		ID:       u.ID,
		Email:    u.Email,
		Password: u.Password,
	}

	userMap[use.ID] = use

	return use, nil
}

// GetByID gets a user by the given ID.
func (db *Database) GetByID(id uint) (*user.User, error) {
	u, ok := userMap[id]
	if !ok {
		return nil, user.ErrUserNotFound
	}

	return u, nil
}

// GetByEmail gets a user by the given email.
func (db *Database) GetByEmail(email string) (*user.User, error) {
	// Loop through users.
	for _, v := range userMap {
		if v.Email == email {
			return v, nil
		}
	}

	return nil, user.ErrUserNotFound
}

// Update updates an invoice.
func (db *Database) Update(u *user.User) (*user.User, error) {
	userMap[u.ID] = u

	return u, nil
}
