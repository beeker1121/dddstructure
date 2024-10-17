package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	apictx "dddstructure/cmd/api/context"
	v1 "dddstructure/cmd/api/v1"
	"dddstructure/dep"
	"dddstructure/service"
	"dddstructure/service/interfaces"
	"dddstructure/storage/mysql"

	"github.com/beeker1121/httprouter"
)

func main() {
	fmt.Println("running...")

	// Create a new MySQL storage implementation.
	fmt.Println("[+] Creating new MySQL storage implementation...")
	store := mysql.New(&sql.DB{})

	// Create a new service.
	fmt.Println("[+] Creating new service...")
	serv := service.New(store)

	// Create services interface.
	servi := interfaces.NewService(interfaces.NewServiceParams{
		Merchant: serv.Merchant,
		Invoice:  serv.Invoice,
	})

	// Set services interface to all services.
	serv.SetServices(servi)

	// Register dependencies.
	dep.RegisterMerchant(serv.Merchant)
	dep.RegisterUser(serv.User)
	dep.RegisterInvoice(serv.Invoice)
	dep.RegisterTransaction(serv.Transaction)

	// Create a new router.
	router := httprouter.New()

	// Create a new API context.
	ac := apictx.New(serv)

	// Create a new v1 API.
	v1.New(ac, router)

	// Create a new HTTP server.
	server := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Printf("Running server...")

	// Start the HTTP server.
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
