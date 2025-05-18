package routes

import (
	"github.com/gorilla/mux"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/handlers"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/handlers/auth"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/handlers/permissions"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/handlers/roles"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/handlers/users"
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

	// Permission management routes (admin only)
	permissionRoutes := protectedRoutes.PathPrefix("/permissions").Subrouter()
	permissionRoutes.Use(middleware.RequireAdmin)
	permissionRoutes.HandleFunc("", permissions.ListPermissionsHandler).Methods("GET")
	permissionRoutes.HandleFunc("/{permission_id}", permissions.GetPermissionHandler).Methods("GET")

	// User profile and permissions (authenticated only)
	protectedRoutes.HandleFunc("/me", users.GetProfileHandler).Methods("GET")
	protectedRoutes.HandleFunc("/me/permissions", permissions.GetUserPermissionsHandler).Methods("GET")

	// user management routes
	userRoutes := protectedRoutes.PathPrefix("/users").Subrouter()

	// List users requires admin
	listRoute := userRoutes.PathPrefix("").Subrouter()
	listRoute.Use(middleware.RequireAdmin)
	listRoute.HandleFunc("", users.ListUsersHandler).Methods("GET")

	// Get user details requires admin or self
	detailRoute := userRoutes.PathPrefix("").Subrouter()
	detailRoute.Use(middleware.RequireAdminOrSelf)
	detailRoute.HandleFunc("/{user_id}", users.GetUserHandler).Methods("GET")

	return r
}
