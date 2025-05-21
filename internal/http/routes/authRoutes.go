package routes

import (
	"log"

	"github.com/gorilla/mux"
	"github.com/soyaibzihad10/Developer-Assignment/internal/config"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/handlers/auth"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/middleware"
)

func RegisterAuthRoutes(r *mux.Router) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	authRoutes := r.PathPrefix("/api/v1/auth").Subrouter()
	authRoutes.HandleFunc("/register", auth.RegisterHandler(cfg)).Methods("POST")
	authRoutes.HandleFunc("/login", auth.LoginHandler).Methods("POST")
	authRoutes.HandleFunc("/verify", auth.VerifyEmailHandler).Methods("GET")
	authRoutes.HandleFunc("/resend-verification", auth.ResendVerificationHandler(cfg)).Methods("POST")
	authRoutes.HandleFunc("/password-reset-request", auth.RequestPasswordResetHandler(cfg)).Methods("POST")
	authRoutes.HandleFunc("/password-reset", auth.ResetPasswordHandler).Methods("POST")

	protected := r.PathPrefix("/api/v1/auth").Subrouter()
	protected.Use(middleware.IsAuthenticated)
	protected.HandleFunc("/check", auth.CheckAuthHandler).Methods("GET")
	protected.HandleFunc("/logout", auth.LogoutHandler).Methods("POST")
}
