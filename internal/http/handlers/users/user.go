package users

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/soyaibzihad10/Developer-Assignment/internal/database"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/middleware"
)

// ListUsersHandler handles GET /api/v1/users
func ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := database.ListUsers()
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"users":  users,
		"total":  len(users),
	})
}

// GetUserHandler handles GET /api/v1/users/{user_id}
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	// Check if user is requesting their own data or has admin/moderator access
	requestorID := r.Context().Value(middleware.UserIDKey).(string)
	requestorType := r.Context().Value(middleware.UserTypeKey).(string)

	if userID != requestorID && requestorType == "user" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	user, err := database.GetUserByID(database.DB, userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"user":   user,
	})
}
