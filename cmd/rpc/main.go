package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	rpcctx "dddstructure/cmd/rpc/context"
	"dddstructure/cmd/rpc/handlers/transaction"
	"dddstructure/dep"
	"dddstructure/service"
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

	// Register dependencies.
	dep.RegisterMerchant(serv.Merchant)
	dep.RegisterUser(serv.User)
	dep.RegisterInvoice(serv.Invoice)
	dep.RegisterProcessor(serv.Processor)
	dep.RegisterTransaction(serv.Transaction)

	// Create a new router.
	router := httprouter.New()

	// Create a new API context.
	rc := rpcctx.New(serv)

	// Create the RPC handlers.
	transaction.New(rc, router)

	// Create a new HTTP server.
	server := &http.Server{
		Addr:           ":8083",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("Running server...")

	// Start the HTTP server.
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
