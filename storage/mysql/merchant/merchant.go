package merchant

import (
	"database/sql"
	"errors"
	"fmt"

	"dddstructure/storage/merchant"
)

// merchantMap acts as a mock MySQL database for merchants.
var merchantMap map[uint]*merchant.Merchant = make(map[uint]*merchant.Merchant)

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

// Create creates a new merchant.
func (db *Database) Create(params *merchant.CreateParams) (*merchant.Merchant, error) {
	m := &merchant.Merchant{
		ID:    params.ID,
		Name:  params.Name,
		Email: params.Email,
	}

	merchantMap[m.ID] = m

	fmt.Println("Created merchant and added to MySQL database...")

	return m, nil
}

// GetByID gets a merchant by the given ID.
func (db *Database) GetByID(id uint) (*merchant.Merchant, error) {
	m, ok := merchantMap[id]
	if !ok {
		return nil, errors.New("could not find merchant")
	}

	fmt.Println("Got merchant from MySQL database...")

	return m, nil
}
