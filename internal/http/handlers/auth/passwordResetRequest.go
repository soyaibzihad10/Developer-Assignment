package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/soyaibzihad10/Developer-Assignment/internal/config"
	"github.com/soyaibzihad10/Developer-Assignment/internal/database"
	"github.com/soyaibzihad10/Developer-Assignment/internal/email"
)

type PasswordResetRequest struct {
	Email string `json:"email"`
}

// RequestPasswordResetHandler creates a handler that uses the provided Config
func RequestPasswordResetHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PasswordResetRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		user, err := database.GetUserByEmail(database.DB, req.Email)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "If your email is registered, you will receive a password reset link",
			})
			return
		}

		resetToken := database.GenerateRandomToken()
		expiry := time.Now().Add(15 * time.Minute)

		err = database.SavePasswordResetToken(database.DB, user.ID, resetToken, expiry)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		resetLink := fmt.Sprintf("%s/reset-password?token=%s", cfg.App.AppURL, resetToken)

		go func() {
			err = email.SendPasswordResetEmail(cfg, user.Email, resetLink)
			if err != nil {
				http.Error(w, "Could not send reset email", http.StatusInternalServerError)
			}
		}()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "If your email is registered, you will receive a password reset link",
		})
	}
}
