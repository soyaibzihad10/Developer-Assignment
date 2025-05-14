package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

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

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
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

	// Call CreateUser to insert into DB
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

	err = email.SendVerificationEmail(user.Email, user.VerificationToken)
	if err != nil {
		http.Error(w, "Could not send verification email", http.StatusInternalServerError)
		return
	}

	// Successful response
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      "success",
		"status_code": 200,
		"user": map[string]interface{}{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
		},
	})
}
