package database

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"github.com/soyaibzihad10/Developer-Assignment/internal/models"
)

var ErrUserNotFound = errors.New("user not found")

func GenerateVerificationToken() (string, error) {
	token := make([]byte, 16) // 16-byte token
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}

func UpdateUserVerificationToken(db *sql.DB, user *models.User) error {
    query := "UPDATE users SET verification_token = $1, token_expiry = $2 WHERE id = $3"
    _, err := db.Exec(query, user.VerificationToken, user.TokenExpiry, user.ID)
    return err
}

func CreateUser(db *sql.DB, user *models.User) error {
	query := `INSERT INTO users (
        username, email, first_name, last_name, password_hash,
        email_verified, user_type, active, verification_token, token_expiry, created_at, updated_at
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
    RETURNING id`

	now := time.Now()

	// Generate verification token
	token, err := GenerateVerificationToken()
	if err != nil {
		return err
	}

	// Default/auto fields set
	user.CreatedAt = now
	user.UpdatedAt = now
	user.EmailVerified = false // new users are not verified by default
	user.UserType = "user"     // default role
	user.Active = true         // default active user
	user.VerificationToken = token
	user.TokenExpiry = now.Add(24 * time.Hour) // Token valid for 24 hours

	err = db.QueryRow(query,
		user.Username,
		user.Email,
		user.FirstName,
		user.LastName,
		user.PasswordHash, // hashed password
		user.EmailVerified,
		user.UserType,
		user.Active,
		user.VerificationToken,
		user.TokenExpiry,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	return err
}

func GetUserByEmail(db *sql.DB, email string) (*models.User, error) {
	user := &models.User{}
	query := "SELECT id, username, first_name, last_name, email, password_hash, email_verified, user_type, active, created_at, updated_at FROM users WHERE email = $1"
	err := db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
		&user.EmailVerified,
		&user.UserType,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}
