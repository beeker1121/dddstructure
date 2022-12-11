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
func (db *Database) Create(m *merchant.Merchant) (*merchant.Merchant, error) {
	merch := &merchant.Merchant{
		ID:    m.ID,
		Name:  m.Name,
		Email: m.Email,
	}

	merchantMap[merch.ID] = merch

	fmt.Println("Created merchant and added to MySQL database...")

	return merch, nil
}

// GetByID gets an merchant by the given ID.
func (db *Database) GetByID(id uint) (*merchant.Merchant, error) {
	m, ok := merchantMap[id]
	if !ok {
		return nil, errors.New("could not find merchant")
	}

	fmt.Println("Got merchant from MySQL database...")

	return m, nil
}
