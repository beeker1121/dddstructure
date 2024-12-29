package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"dddstructure/cmd/api/config"
	apictx "dddstructure/cmd/api/context"
	v1 "dddstructure/cmd/api/v1"
	"dddstructure/service"
	storagemysql "dddstructure/storage/mysql"

	"github.com/beeker1121/httprouter"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Parse the API configuration file.
	cfg, err := config.ParseConfigFile("config.json")
	if err != nil {
		panic(err)
	}

	// Get the configuration environment variables.
	cfg.DBHost = os.Getenv("DB_HOST")
	cfg.DBPort = os.Getenv("DB_PORT")
	cfg.DBName = os.Getenv("DB_NAME")
	cfg.DBUser = os.Getenv("DB_USER")
	cfg.DBPass = os.Getenv("DB_PASS")
	cfg.APIHost = os.Getenv("API_HOST")
	cfg.APIPort = os.Getenv("API_PORT")
	cfg.JWTSecret = os.Getenv("JWT_SECRET")

	/* // Create a new mock storage implementation.
	fmt.Println("[+] Creating new mock storage implementation...")
	store := mock.New(&sql.DB{}) */

	// Connect to the MySQL database.
	db, err := sql.Open("mysql", cfg.DBUser+":"+cfg.DBPass+"@tcp("+cfg.DBHost+":"+cfg.DBPort+")/"+cfg.DBName+"?parseTime=true")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Test database connection.
	if err := db.Ping(); err != nil {
		panic(err)
	}

	// Create a new MySQL storage implementation.
	store := storagemysql.New(db)

	// Create a new service.
	fmt.Println("[+] Creating new service...")
	serv := service.New(store)

	// Create a new router.
	router := httprouter.New()

	// Create a new API context.
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	ac := apictx.New(cfg, logger, serv)

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

	fmt.Printf("[+] Running server on port 8080...\n")

	// Start the HTTP server.
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
