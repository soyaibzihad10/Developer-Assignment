package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/soyaibzihad10/Developer-Assignment/internal/database"
)

func VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	// Debug logging
	log.Printf("Received verification request with token length: %d", len(token))

	// Query to get user details first
	var userID string
	var email string
	userQuery := `
        SELECT id, email 
        FROM users 
        WHERE verification_token = $1 
        AND token_expiry > NOW()
        AND email_verified = false`

	err := database.DB.QueryRow(userQuery, token).Scan(&userID, &email)
	if err != nil {
		log.Printf("Error finding user with token: %v", err)
		http.Error(w, "Invalid or expired token", http.StatusBadRequest)
		return
	}

	// Update user verification status
	updateQuery := `
        UPDATE users 
        SET 
            email_verified = true,
            verification_token = NULL,
            token_expiry = NULL
        WHERE id = $1`

	result, err := database.DB.Exec(updateQuery, userID)
	if err != nil {
		log.Printf("Error updating user verification status: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("No user updated for ID: %s", userID)
		http.Error(w, "Verification failed", http.StatusBadRequest)
		return
	}

	log.Printf("Successfully verified email for user: %s", email)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Email verified successfully",
	})
}
