package main

import (
	"log"
	"net/http"

	"department-management/api"
	"department-management/config"
	"department-management/db"
	"github.com/gorilla/handlers"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the database connection
	database, err := db.InitDB(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// Initialize the router
	router := api.RegisterRoutes()

	// Start the server
	log.Printf("Starting server on port %s...", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, handlers.CORS()(router)))
}
