package routes

import (
	"github.com/gorilla/mux"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/handlers"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/handlers/auth"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/handlers/roles"
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

	protectedRoutes := r.PathPrefix("/api/v1").Subrouter() // Fixed typo here
	protectedRoutes.Use(middleware.IsAuthenticated)
	protectedRoutes.HandleFunc("/auth/check", auth.CheckAuthHandler).Methods("GET")
	protectedRoutes.HandleFunc("/auth/logout", auth.LogoutHandler).Methods("POST")

	// Role management routes
	roleRoutes := protectedRoutes.PathPrefix("/roles").Subrouter()
	roleRoutes.Use(middleware.RequireAdmin)

	roleRoutes.HandleFunc("", roles.ListRolesHandler).Methods("GET")
	roleRoutes.HandleFunc("/{role_id}", roles.GetRoleHandler).Methods("GET")
	roleRoutes.HandleFunc("", roles.CreateRoleHandler).Methods("POST")
	roleRoutes.HandleFunc("/{role_id}", roles.UpdateRoleHandler).Methods("PUT")
	roleRoutes.HandleFunc("/{role_id}", roles.DeleteRoleHandler).Methods("DELETE")

	return r
}
