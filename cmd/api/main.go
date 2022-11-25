package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	apictx "dddstructure/cmd/api/context"
	v1 "dddstructure/cmd/api/v1"
	"dddstructure/service"
	"dddstructure/service/core"
	"dddstructure/storage/mysql"

	"github.com/beeker1121/httprouter"
)

func main() {
	// Create a new mysql storage implementation.
	fmt.Println("[+] Creating new MySQL storage implementation...")
	store := mysql.New(&sql.DB{})

	// Create a new core service.
	fmt.Println("[+] Creating new core service...")
	coreserv := core.New(store)

	// Create a new service.
	fmt.Println("[+] Creating new service...")
	serv := service.New(coreserv)

	// Create a new router.
	router := httprouter.New()

	// Create a new API context.
	ac := apictx.New(serv)

	// Create a new API v1.
	v1.New(ac, router)

	// Create a new HTTP server.
	server := &http.Server{
		Addr:           ":9280",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("[+] Running server on :9280...")

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
