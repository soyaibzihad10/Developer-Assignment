package routes

import (
	"github.com/gorilla/mux"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/handlers"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/handlers/auth"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// Health checker route
	r.HandleFunc("/ping", handlers.PingHandler).Methods("GET")

	// auth routes
	authRoutes := r.PathPrefix("/api/v1/auth").Subrouter()
	authRoutes.HandleFunc("/register", auth.RegisterHandler).Methods("POST")
	authRoutes.HandleFunc("/login", auth.LoginHandler).Methods("POST")
	return r
}
