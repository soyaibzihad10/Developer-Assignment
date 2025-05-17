package roles

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/soyaibzihad10/Developer-Assignment/internal/database"
	"github.com/soyaibzihad10/Developer-Assignment/internal/models"
)

func ListRolesHandler(w http.ResponseWriter, r *http.Request) {
	roles, err := database.ListRoles()
	if err != nil {
		http.Error(w, "Failed to fetch roles", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.RoleList{
		Roles: roles,
		Total: len(roles),
	})
}

func GetRoleHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["role_id"]
	role, err := database.GetRoleByID(id)
	if err != nil {
		http.Error(w, "Role not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(role)
}

func CreateRoleHandler(w http.ResponseWriter, r *http.Request) {
	var role models.Role
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if role.Name == "" {
		http.Error(w, "Role name is required", http.StatusBadRequest)
		return
	}

	if err := database.CreateRole(&role); err != nil {
		http.Error(w, "Failed to create role", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(role)
}

func UpdateRoleHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["role_id"]
	var role models.Role
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := database.UpdateRole(id, &role); err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "Role not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update role", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(role)
}

func DeleteRoleHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["role_id"]
	if err := database.DeleteRole(id); err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "Role not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete role", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
