package main

import (
	"log"
	"net/http"

	"github.com/soyaibzihad10/Developer-Assignment/internal/config"
	"github.com/soyaibzihad10/Developer-Assignment/internal/database"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/routes"
)

var cnf *config.Config

func init() {
	var err error
	cnf, err = config.LoadConfig()
	if err != nil {
		log.Println("Config func does not working well")
	}

	if err := database.ConnDB(cnf.Database); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	database.CreateSystemAdminIfNotExists(cnf.Admin)

	if err := database.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
}

func main() {
	log.Println("Server initialized")
	defer database.DB.Close()

	router := routes.SetupRoutes()

	corsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		router.ServeHTTP(w, r)
	})


	port := ":8080"
	log.Println("Server is running on http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, corsHandler))
}
