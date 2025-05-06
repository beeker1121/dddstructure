package user

import (
	"context"
	"database/sql"

	"dddstructure/storage/mysql/models"
	"dddstructure/storage/user"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

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
	// Map to model.
	model := models.User{
		ID:       u.ID,
		Email:    u.Email,
		Password: u.Password,
	}

	// Insert into database.
	err := model.Insert(context.Background(), db.db, boil.Infer())
	if err != nil {
		return nil, err
	}

	return u, nil
}

// GetByID gets a user by the given ID.
func (db *Database) GetByID(id uint) (*user.User, error) {
	modelu, err := models.Users(qm.Where("id=?", id)).One(context.Background(), db.db)
	if err == sql.ErrNoRows {
		return nil, user.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	// Map to user type.
	u := &user.User{
		ID:       modelu.ID,
		Email:    modelu.Email,
		Password: modelu.Password,
	}

	return u, nil
}

// GetByEmail gets a user by the given email.
func (db *Database) GetByEmail(email string) (*user.User, error) {
	modelu, err := models.Users(qm.Where("email=?", email)).One(context.Background(), db.db)
	if err == sql.ErrNoRows {
		return nil, user.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	// Map to user type.
	u := &user.User{
		ID:       modelu.ID,
		Email:    modelu.Email,
		Password: modelu.Password,
	}

	return u, nil
}

// Update updates an invoice.
func (db *Database) Update(u *user.User) (*user.User, error) {
	// Map to model.
	model := models.User{
		ID:       u.ID,
		Email:    u.Email,
		Password: u.Password,
	}

	// Update in database.
	_, err := model.Update(context.Background(), db.db, boil.Infer())
	if err != nil {
		return nil, err
	}

	return u, nil
}
