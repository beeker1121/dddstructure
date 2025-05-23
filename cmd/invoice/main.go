package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"dddstructure/proto"
	"dddstructure/service"
	"dddstructure/storage/mock"
)

func main() {
	fmt.Println("running...")

	// Create a new logger.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Create a new mock storage implementation.
	fmt.Println("[+] Creating new mock storage implementation...")
	store := mock.New(&sql.DB{})

	// Create a new service.
	fmt.Println("[+] Creating new service...")
	serv := service.New(store, logger)

	// Create a user.
	u, err := serv.User.Create(&proto.UserCreateParams{
		Email:    "johndoe@gmail.com",
		Password: "TestPassword123",
	})
	if err != nil {
		panic(err)
	}

	// Create an invoice.
	i, err := serv.Invoice.Create(&proto.InvoiceCreateParams{
		UserID:         u.ID,
		PaymentMethods: []proto.InvoicePaymentMethod{proto.InvoicePaymentMethodCard},
		BillTo: proto.InvoiceBillTo{
			FirstName: "Bill",
			LastName:  "Smith",
		},
		PayTo: proto.InvoicePayTo{
			FirstName: "John",
			LastName:  "Doe",
		},
		LineItems: []proto.InvoiceLineItem{
			{
				Quantity: 1,
				Price:    100,
			},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("[+] New invoice: %+v\n", *i)

	// Pay an invoice, will call transaction.Process service.
	i, err = serv.Invoice.Pay(i.ID, &proto.InvoicePayParams{
		Amount: 100,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("[+] Paid invoice: %+v\n", *i)

	// Process a transaction, will call invoice.Update service.
	t, err := serv.Transaction.Process(&proto.TransactionProcessParams{
		UserID:    i.UserID,
		Type:      "refund",
		Amount:    100,
		InvoiceID: i.ID,
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
