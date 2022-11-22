package main

import (
	"database/sql"
	"fmt"

	"dddstructure/service"
	"dddstructure/service/accounting"
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

	/* --- Create a new accounting entry --- */

	// Get the user for this merchant.
	fmt.Println("[+] Getting user for merchant...")
	u, err := serv.User.GetByID(1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[+] Got user: %+v\n", *u)

	// Create a new accounting entry.
	fmt.Printf("[+] Creating a new accounting entry for merchant '%v' and user '%v'...\n", m.ID, u.ID)
	a, err := serv.Accounting.Create(&accounting.CreateParams{
		ID:         1,
		MerchantID: m.ID,
		UserID:     u.ID,
		AmountDue:  100,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("[+] Created accounting entry: %+v\n", *a)

	/* --- Handle billing --- */

	// Get billing information for all merchants.
	fmt.Println("[+] Getting billing information for all merchants...")
	b, err := serv.Billing.GetMerchantAmountsDue()
	if err != nil {
		panic(err)
	}
	for _, v := range b {
		fmt.Printf("[+] Got billing information for merchant: %+v\n", *v)
	}

	// // Mark as paid.
	// fmt.Println("[+] Marking as paid...")
	// err = serv.Billing.AddPayment(billing.AddPaymentParams{
	// 	MerchantID: m.ID,
	// 	UserID:     u.ID,
	// 	Amount:     100,
	// })

	// // Get billing information for all merchants.
	// fmt.Println("[+] Getting billing information for all merchants...")
	// b, err = serv.Billing.GetMerchantAmountsDue()
	// if err != nil {
	// 	panic(err)
	// }
	// for _, v := range b {
	// 	fmt.Printf("[+] Got billing information for merchant: %+v\n", *v)
	// }
}
