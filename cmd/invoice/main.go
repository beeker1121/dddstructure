package main

import (
	"database/sql"
	"fmt"

	"dddstructure/service"
	"dddstructure/storage/mysql"
)

func main() {
	fmt.Println("running...")

	/* --- Setup --- */

	// Create a new mysql storage implementation.
	fmt.Println("[+] Creating new MySQL storage implementation...")
	store := mysql.New(&sql.DB{})

	// Create a new service.
	fmt.Println("[+] Creating new service...")
	serv := service.New(store)

	/* --- Create merchant --- */
	serv.Invoice.GetByID(1)
}
