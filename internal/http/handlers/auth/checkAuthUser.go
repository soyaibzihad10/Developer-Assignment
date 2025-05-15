package auth

import (
	"encoding/json"
	"net/http"
)

func CheckAuthHandler(w http.ResponseWriter, r *http.Request) {
	// info from the request header set by the middleware
	userID := r.Header.Get("user_id")
	userEmail := r.Header.Get("user_email")

	// if userID is empty, it means the user is not authenticated
	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":        "error",
			"authenticated": false,
			"message":       "User not authenticated",
		})
		return
	}

	// successful authentication response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":        "success",
		"authenticated": true,
		"user": map[string]string{
			"id":    userID,
			"email": userEmail,
		},
	})
}
