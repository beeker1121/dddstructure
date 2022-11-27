package main

import (
	"database/sql"
	"fmt"

	"dddstructure/service"
	"dddstructure/service/core"
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

	/* --- Create merchant --- */
	serv.Merchant.GetByID(1)
}
