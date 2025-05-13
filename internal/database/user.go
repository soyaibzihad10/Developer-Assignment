package database

import (
	"database/sql"
	"time"

	"github.com/soyaibzihad10/Developer-Assignment/internal/models"
)

func CreateUser(db *sql.DB, user *models.User) error {
	query := `INSERT INTO users (
		username, email, first_name, last_name, password_hash,
		email_verified, user_type, active, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id`

	now := time.Now()

	// Default/auto fields set
	user.CreatedAt = now
	user.UpdatedAt = now
	user.EmailVerified = false   // new users are not verified by default
	user.UserType = "user"       // default role
	user.Active = true           // default active user

	err := db.QueryRow(query,
		user.Username,
		user.Email,
		user.FirstName,
		user.LastName,
		user.PasswordHash,     // hashed password
		user.EmailVerified,
		user.UserType,
		user.Active,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	return err
}
