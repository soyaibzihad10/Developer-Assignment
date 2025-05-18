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

// password reset related
func GenerateRandomToken() string {
    b := make([]byte, 32)
    rand.Read(b)
    return hex.EncodeToString(b)
}

func SavePasswordResetToken(db *sql.DB, userID string, token string, expiry time.Time) error {
    query := `UPDATE users SET 
        reset_token = $1,
        reset_token_expiry = $2
        WHERE id = $3`
    _, err := db.Exec(query, token, expiry, userID)
    return err
}

func ValidatePasswordResetToken(db *sql.DB, token string) (string, error) {
    var userID string
    query := `SELECT id FROM users 
        WHERE reset_token = $1 
        AND reset_token_expiry > NOW()`
    err := db.QueryRow(query, token).Scan(&userID)
    return userID, err
}

func UpdateUserPassword(db *sql.DB, userID string, hashedPassword string) error {
    query := `UPDATE users SET 
        password_hash = $1,
        reset_token = NULL,
        reset_token_expiry = NULL
        WHERE id = $2`
    _, err := db.Exec(query, hashedPassword, userID)
    return err
}


// varification related
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
    // Start a transaction
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // Insert user
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
    user.EmailVerified = false
    user.UserType = "user"
    user.Active = true
    user.VerificationToken = token
    user.TokenExpiry = now.Add(24 * time.Hour)

    err = tx.QueryRow(query,
        user.Username,
        user.Email,
        user.FirstName,
        user.LastName,
        user.PasswordHash,
        user.EmailVerified,
        user.UserType,
        user.Active,
        user.VerificationToken,
        user.TokenExpiry,
        user.CreatedAt,
        user.UpdatedAt,
    ).Scan(&user.ID)

    if err != nil {
        return err
    }

    // Assign default user role
    roleQuery := `
        INSERT INTO user_roles (user_id, role_id, assigned_by)
        SELECT $1, id, $1
        FROM roles
        WHERE name = 'user'
    `
    _, err = tx.Exec(roleQuery, user.ID)
    if err != nil {
        return err
    }

    // Commit transaction
    return tx.Commit()
}

func GetUserByID(db *sql.DB, id string) (*models.User, error) {
    user := &models.User{}
    err := db.QueryRow(`
        SELECT id, username, email, first_name, last_name, 
               user_type, email_verified, active, created_at, updated_at
        FROM users 
        WHERE id = $1 AND deletion_requested = false
    `, id).Scan(
        &user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName,
        &user.UserType, &user.EmailVerified, &user.Active, &user.CreatedAt, &user.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, ErrNotFound
    }
    if err != nil {
        return nil, err
    }

    return user, nil
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
