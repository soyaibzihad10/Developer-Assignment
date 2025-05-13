package routes

import (
	"github.com/gorilla/mux"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/handlers"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// Health check route
	r.HandleFunc("/ping", handlers.PingHandler).Methods("GET")
	return r
}
