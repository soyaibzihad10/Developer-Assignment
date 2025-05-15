package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/soyaibzihad10/Developer-Assignment/internal/config"
	"github.com/soyaibzihad10/Developer-Assignment/internal/database"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/routes"
)

var cnf *config.Config

func init() {
	var err error
	cnf, err = config.LoadConfig() // get .env variables
	if err != nil {
		log.Println("Config func does not working well")
	}

	// connect to database
	if err := database.ConnDB(cnf.Database); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// run migrations
	if err := database.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
}

func main() {
	fmt.Println("AffPilot Auth Service starting...")
	log.Println("Server initialized")

	defer database.DB.Close()

	r := routes.SetupRoutes()

	// Start server
	port := ":8080"
	log.Println("Server is running on http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, r))
}
