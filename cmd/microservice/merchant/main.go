package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	msctx "dddstructure/cmd/microservice/merchant/context"
	"dddstructure/cmd/microservice/merchant/rest"
	"dddstructure/cmd/microservice/merchant/service"
	"dddstructure/cmd/microservice/merchant/storage"
	"dddstructure/dep"
	"dddstructure/storage/mysql/merchant"

	"github.com/beeker1121/httprouter"
)

func main() {
	fmt.Println("running...")

	// Create a new MySQL database.
	db := &sql.DB{}

	// Create a new user MySQL storage implementation.
	merchantdb := merchant.New(db)

	// Create a new storage implementation.
	store := &storage.Storage{
		Merchant: merchantdb,
	}

	// Create a new service.
	serv := service.New(store)

	dep.RegisterUser(serv.User)

	// Create a new router.
	router := httprouter.New()

	// Create a new microservice context.
	mc := msctx.New(serv)

	// Create a new v1 API.
	rest.New(mc, router)

	// Create a new HTTP server.
	server := &http.Server{
		Addr:           ":8082",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("Running merchant microservice server...")

	// Start the HTTP server.
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
