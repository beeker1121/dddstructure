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

	/* --- Create and get merchant --- */

	// Create a new merchant.
	fmt.Println("[+] Creating a merchant via the service...")
	m, err := serv.Merchant.Create(&merchant.CreateParams{
		ID:    1,
		Name:  "John Doe",
		Email: "johndoe@gmail.com",
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("[+] Created merchant: %+v\n", *m)

	// Get the merchant.
	fmt.Println("[+] Getting merchant via the service...")
	m, err = serv.Merchant.GetByID(1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[+] Got merchant: %+v\n", *m)

	/* --- Handle billing --- */

	// Get the user for this merchant.
	fmt.Println("[+] Getting user for merchant...")
	u, err := serv.User.GetByID(1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[+] Got user: %+v\n", *u)

	// Add amount owed via billing service.
	fmt.Println("[+] Adding amount owed via billing service...")
	if err := serv.Billing.AddAmountDue(m.ID, u.ID, 100); err != nil {
		panic(err)
	}

	// Get billing information for all merchants.
	fmt.Println("[+] Getting billing information for all merchants...")
	b, err := serv.Billing.GetMerchantAmountsDue()
	if err != nil {
		panic(err)
	}
	for _, v := range b {
		fmt.Printf("[+] Got billing information for merchant: %+v\n", *v)

		// Add amount paid.
		fmt.Println("[+] Adding 100 as amount paid...")
		err = serv.Billing.AddAmountPaid(v.ID, 100)
		if err != nil {
			panic(err)
		}
	}

	// Getting accounting entry for the merchant we created.
	fmt.Printf("[+] Getting accounting entry for merchant '%v'...\n", m.ID)
	a, err := serv.Accounting.GetByID(1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[+] Got accounting entry: %+v\n", *a)
}
