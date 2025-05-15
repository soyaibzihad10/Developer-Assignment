package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/soyaibzihad10/Developer-Assignment/internal/database"
	"github.com/soyaibzihad10/Developer-Assignment/internal/email"
)

type PasswordResetRequest struct {
	Email string `json:"email"`
}

func RequestPasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	var req PasswordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// ইউজার চেক
	user, err := database.GetUserByEmail(database.DB, req.Email)
	if err != nil {
		// সিকিউরিটির জন্য একই মেসেজ দিন
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "If your email is registered, you will receive a password reset link",
		})
		return
	}

	// রিসেট টোকেন জেনারেট
	resetToken := database.GenerateRandomToken()
	expiry := time.Now().Add(15 * time.Minute)

	// টোকেন ডাটাবেসে সেভ
	err = database.SavePasswordResetToken(database.DB, user.ID, resetToken, expiry)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// ইমেইল পাঠান
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", os.Getenv("APP_URL"), resetToken)
	err = email.SendPasswordResetEmail(user.Email, resetLink)
	if err != nil {
		http.Error(w, "Could not send reset email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "If your email is registered, you will receive a password reset link",
	})
}
