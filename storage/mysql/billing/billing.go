package billing

import (
	"database/sql"
	"errors"
	"fmt"

	"dddstructure/storage/billing"
)

var billingMap = map[uint]*billing.MerchantAmountsDue{1: &billing.MerchantAmountsDue{
	ID:           1,
	MerchantID:   1,
	MerchantName: "John Doe",
	AmountDue:    100,
}}

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

// GetMerchantAmountsDue gets all merchant amounts due.
func (db *Database) GetMerchantAmountsDue() ([]*billing.MerchantAmountsDue, error) {
	// Would call custom SQL query here to JOIN accounting, merchant, and user
	// tables. The only entry in billingMap mocks this.
	//
	// The idea here is that the storage level domain only cares about getting
	// and setting data to and from the database. If we need to create a
	// separate storage implementation that does a JOIN on multiple tables, ie
	// just SELECT and JOIN... this does not implement any business logic and
	// is deemed ok to do.
	m, ok := billingMap[1]
	if !ok {
		return nil, errors.New("could not find billing entry")
	}

	fmt.Println("Got billing entry from MySQL database...")

	return []*billing.MerchantAmountsDue{m}, nil
}
