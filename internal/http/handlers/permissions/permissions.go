package permissions

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/soyaibzihad10/Developer-Assignment/internal/database"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/middleware"
)

// ListPermissionsHandler handles GET /api/v1/permissions
func ListPermissionsHandler(w http.ResponseWriter, r *http.Request) {
	permissions, err := database.ListPermissions()
	if err != nil {
		http.Error(w, "Failed to fetch permissions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      "success",
		"permissions": permissions,
		"total":       len(permissions),
	})
}

// GetPermissionHandler handles GET /api/v1/permissions/{permission_id}
func GetPermissionHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["permission_id"]

	permission, err := database.GetPermissionByID(id)
	if err != nil {
		http.Error(w, "Permission not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "success",
		"permission": permission,
	})
}

// GetUserPermissionsHandler handles GET /api/v1/me/permissions
func GetUserPermissionsHandler(w http.ResponseWriter, r *http.Request) {
	userType := r.Context().Value(middleware.UserTypeKey).(string)
	if userType == "" {
		http.Error(w, "User type not found", http.StatusUnauthorized)
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)
	permissions, err := database.GetUserPermissions(userID)
	if err != nil {
		http.Error(w, "Failed to fetch user permissions", http.StatusInternalServerError)
		return
	}

	fmt.Printf("chill %s", userID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      "success",
		"permissions": permissions,
		"total":       len(permissions),
	})
}
