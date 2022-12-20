package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	msctx "dddstructure/cmd/microservice/user/context"
	"dddstructure/cmd/microservice/user/rest"
	"dddstructure/cmd/microservice/user/service"
	"dddstructure/cmd/microservice/user/storage"
	"dddstructure/storage/mysql/user"

	"github.com/beeker1121/httprouter"
)

func main() {
	fmt.Println("running...")

	// Create a new MySQL database.
	db := &sql.DB{}

	// Create a new user MySQL storage implementation.
	userdb := user.New(db)

	// Create a new storage implementation.
	store := &storage.Storage{
		User: userdb,
	}

	// Create a new service.
	serv := service.New(store)

	// Create a new router.
	router := httprouter.New()

	// Create a new microservice context.
	mc := msctx.New(serv)

	// Create a new v1 API.
	rest.New(mc, router)

	// Create a new HTTP server.
	server := &http.Server{
		Addr:           ":8081",
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
