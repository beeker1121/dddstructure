package main

import (
	"database/sql"
	"fmt"

	"dddstructure/dep"
	"dddstructure/proto"
	"dddstructure/service"
	"dddstructure/storage/mysql"
)

func main() {
	fmt.Println("running...")

	// Create a new MySQL storage implementation.
	fmt.Println("[+] Creating new MySQL storage implementation...")
	store := mysql.New(&sql.DB{})

	// Create a new service.
	fmt.Println("[+] Creating new service...")
	serv := service.New(store)

	// Register dependencies.
	dep.RegisterMerchant(serv.Merchant)
	dep.RegisterUser(serv.User)
	dep.RegisterInvoice(serv.Invoice)
	dep.RegisterProcessor(serv.Processor)
	dep.RegisterTransaction(serv.Transaction)

	// Create a merchant.
	m, err := serv.Merchant.Create(&proto.Merchant{
		Name:  "John Doe",
		Email: "johndoe@gmail.com",
	})

	// Create an invoice.
	i, err := serv.Invoice.Create(&proto.Invoice{
		MerchantID: m.ID,
		BillTo:     "Bill Smith",
		PayTo:      m.Name,
		AmountDue:  100,
		AmountPaid: 0,
		Status:     "pending",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("[+] New invoice: %+v\n", *i)

	// Pay an invoice, will call transaction.Process service.
	i, err = serv.Invoice.Pay(i.ID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[+] Paid invoice: %+v\n", *i)

	// Process a transaction, will call invoice.Update service.
	t, err := serv.Transaction.Process(&proto.Transaction{
		MerchantID:     i.MerchantID,
		Type:           "refund",
		CardType:       "visa",
		AmountCaptured: 100,
		InvoiceID:      i.ID,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("[+] New transaction processed: %+v\n", *t)

	// Get the invoice again.
	i, err = serv.Invoice.GetByID(i.ID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[+] Invoice after transaction refund: %+v\n", *i)
}
