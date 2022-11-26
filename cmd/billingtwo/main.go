package main

import (
	"database/sql"
	"fmt"

	"dddstructure/service"
	"dddstructure/service/core"
	"dddstructure/service/merchant"
	"dddstructure/storage/mysql"
)

func main() {
	fmt.Println("running...")

	/* --- Setup --- */

	// Create a new mysql storage implementation.
	fmt.Println("[+] Creating new MySQL storage implementation...")
	store := mysql.New(&sql.DB{})

	// Create a new core service.
	fmt.Println("[+] Creating new core service...")
	coreserv := core.New(store)

	// Create a new service.
	fmt.Println("[+] Creating new service...")
	serv := service.New(coreserv)

	/* --- Create merchant and add amount due --- */

	// Create new merchant and add amount due.
	if err := createMerchantAndAmountDue(serv); err != nil {
		panic(err)
	}

	/* --- Handle billing --- */

	// Getting accounting entry for the merchant we created.
	fmt.Printf("[+] Getting accounting entry for merchant '%v'...\n", 1)
	a, err := serv.Accounting.GetByID(1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[+] Got accounting entry: %+v\n", *a)

	// Handle billing information for all merchants.
	fmt.Println("[+] Handling billing for all merchants...")
	err = serv.Billing.HandleMerchantBilling()
	if err != nil {
		panic(err)
	}

	// Getting accounting entry for the merchant we created.
	fmt.Printf("[+] Getting accounting entry for merchant '%v'...\n", 1)
	a, err = serv.Accounting.GetByID(1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[+] Got accounting entry: %+v\n", *a)
}

func createMerchantAndAmountDue(s *service.Service) error {
	// Create a new merchant.
	fmt.Println("[+] Creating a merchant via the service...")
	m, err := s.Merchant.Create(&merchant.CreateParams{
		ID:    1,
		Name:  "John Doe",
		Email: "johndoe@gmail.com",
	})
	if err != nil {
		return err
	}

	fmt.Printf("[+] Created merchant: %+v\n", *m)

	// Get the merchant.
	fmt.Println("[+] Getting merchant via the service...")
	m, err = s.Merchant.GetByID(1)
	if err != nil {
		return err
	}
	fmt.Printf("[+] Got merchant: %+v\n", *m)

	// Get the user for this merchant.
	fmt.Println("[+] Getting user for merchant...")
	u, err := s.User.GetByID(1)
	if err != nil {
		return err
	}
	fmt.Printf("[+] Got user: %+v\n", *u)

	// Add amount owed via billing service.
	fmt.Println("[+] Adding amount owed via billing service...")
	if err := s.Billing.AddAmountDue(m.ID, u.ID, 100); err != nil {
		return err
	}

	return nil
}
