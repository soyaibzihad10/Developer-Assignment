package users

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/soyaibzihad10/Developer-Assignment/internal/database"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/middleware"
	"github.com/soyaibzihad10/Developer-Assignment/internal/models"
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

// UpdateUserHandler handles PUT /api/v1/users/{user_id}
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	// Decode request body
	var req models.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update user
	user, err := database.UpdateUser(database.DB, userID, &req)
	if err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"user":   user,
	})
}

// DeleteUserHandler handles DELETE /api/v1/users/{user_id}
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	if err := database.DeleteUser(database.DB, userID); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RequestDeletionHandler handles POST /api/v1/users/{user_id}/request-deletion
func RequestDeletionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	if err := database.RequestUserDeletion(database.DB, userID); err != nil {
		http.Error(w, "Failed to request deletion", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Deletion request submitted successfully",
	})
}

// ChangeRoleHandler handles POST /api/v1/users/{user_id}/role
func ChangeRoleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	var req struct {
		Role string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Println("chill", req.Role)

	if req.Role != "admin" || req.Role != "moderator" || req.Role != "system_admin" {	
		http.Error(w, "Invalid role, you can't change role to admin or system admin or moderator", http.StatusBadRequest)
		return
	}


	if err := database.ChangeUserRole(database.DB, userID, req.Role); err != nil {
		http.Error(w, "Failed to change role", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Role updated successfully",
	})
}

// PromoteToAdminHandler handles POST /api/v1/users/{user_id}/promote/admin
func PromoteToAdminHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	if err := database.PromoteUserToAdmin(database.DB, userID); err != nil {
		http.Error(w, "Failed to promote user to admin", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "User promoted to admin successfully",
	})
}

// PromoteToModeratorHandler handles POST /api/v1/users/{user_id}/promote/moderator
func PromoteToModeratorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	if err := database.PromoteUserToModerator(database.DB, userID); err != nil {
		http.Error(w, "Failed to promote user to moderator", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "User promoted to moderator successfully",
	})
}

// DemoteUserHandler handles POST /api/v1/users/{user_id}/demote
func DemoteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	if err := database.DemoteUser(database.DB, userID); err != nil {
		http.Error(w, "Failed to demote user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "User demoted successfully",
	})
}
