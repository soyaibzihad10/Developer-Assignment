package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/soyaibzihad10/Developer-Assignment/internal/config"
	"github.com/soyaibzihad10/Developer-Assignment/internal/database"
	"github.com/soyaibzihad10/Developer-Assignment/internal/email"
	"github.com/soyaibzihad10/Developer-Assignment/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

// RegisterHandler creates a handler that uses the provided Config
func RegisterHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if req.Email == "" || req.Password == "" || req.Username == "" {
			http.Error(w, "Email, password, and username are required", http.StatusBadRequest)
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Could not hash password", http.StatusInternalServerError)
			return
		}

		user := &models.User{
			Email:        req.Email,
			Username:     req.Username,
			PasswordHash: string(hashedPassword),
		}

		// Call CreateUser to insert into DB and assign role
		err = database.CreateUser(database.DB, user)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value") {
				http.Error(w, "Email or username already exists", http.StatusConflict)
				return
			}
			log.Printf("Error creating user: %v", err)
			http.Error(w, "Could not create user", http.StatusInternalServerError)
			return
		}

		// Send verification email
		go func() {
			err := email.SendVerificationEmail(cfg, user.Email, user.VerificationToken)
			if err != nil {
				log.Printf("Warning: Could not send verification email: %v", err)
			}
		}()

		// Set headers before writing response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		// Return success response
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":      "success",
			"status_code": http.StatusCreated,
			"message":     "User registered successfully. Please check your email for verification.",
			"user": map[string]interface{}{
				"id":       user.ID,
				"email":    user.Email,
				"username": user.Username,
				"active":   user.Active,
			},
		})
	}
}
