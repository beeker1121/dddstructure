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

	// Create a new mysql storage implementation.
	fmt.Println("[+] Creating new MySQL storage implementation...")
	store := mysql.New(&sql.DB{})

	// Create a new core service.
	fmt.Println("[+] Creating new core service...")
	coreserv := core.New(store)

	// Create a new service.
	fmt.Println("[+] Creating new service...")
	serv := service.New(coreserv)

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
}
