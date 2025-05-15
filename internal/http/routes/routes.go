package routes

import (
	"github.com/gorilla/mux"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/handlers"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/handlers/auth"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/middleware"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// Health checker route
	r.HandleFunc("/ping", handlers.PingHandler).Methods("GET")

	// auth routes
	authRoutes := r.PathPrefix("/api/v1/auth").Subrouter()
	authRoutes.HandleFunc("/register", auth.RegisterHandler).Methods("POST")
	authRoutes.HandleFunc("/login", auth.LoginHandler).Methods("POST")
	authRoutes.HandleFunc("/verify", auth.VerifyEmailHandler).Methods("GET")
	authRoutes.HandleFunc("/resend-verification", auth.ResendVerificationHandler).Methods("POST")
	authRoutes.HandleFunc("/password-reset-request", auth.RequestPasswordResetHandler).Methods("POST")
    authRoutes.HandleFunc("/password-reset", auth.ResetPasswordHandler).Methods("POST")
	
	protedtedRoutes := r.PathPrefix("/api/v1").Subrouter()
	protedtedRoutes.Use(middleware.IsAuthenticated)
	protedtedRoutes.HandleFunc("/auth/check", auth.CheckAuthHandler).Methods("GET")
	protedtedRoutes.HandleFunc("/auth/logout", auth.LogoutHandler).Methods("POST")
	
	return r
}
