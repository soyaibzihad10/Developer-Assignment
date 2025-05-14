package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/soyaibzihad10/Developer-Assignment/internal/database"
	"github.com/soyaibzihad10/Developer-Assignment/internal/email"
)

type ResendVerificationRequest struct {
	Email string `json:"email"`
}

func ResendVerificationHandler(w http.ResponseWriter, r *http.Request) {
	var req ResendVerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Fetch user from database
	user, err := database.GetUserByEmail(database.DB, req.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if email is already verified
	if user.EmailVerified {
		http.Error(w, "Email is already verified", http.StatusBadRequest)
		return
	}

	// Generate new verification token
	token, err := database.GenerateVerificationToken()
	if err != nil {
		http.Error(w, "Could not generate verification token", http.StatusInternalServerError)
		return
	}

	// Update user with new token and expiry
	user.VerificationToken = token
	user.TokenExpiry = time.Now().Add(5 * time.Minute) // 5 minutes expiry
	err = database.UpdateUserVerificationToken(database.DB, user)
	if err != nil {
		http.Error(w, "Could not update verification token", http.StatusInternalServerError)
		return
	}

	// Send verification email
	err = email.SendVerificationEmail(user.Email, user.VerificationToken)
	if err != nil {
		http.Error(w, "Could not send verification email", http.StatusInternalServerError)
		return
	}

	// Successful response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Verification email resent successfully"))
}
