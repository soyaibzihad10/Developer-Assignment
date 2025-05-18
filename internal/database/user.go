package database

import (
	"github.com/soyaibzihad10/Developer-Assignment/internal/models"
)

func ListUsers() ([]models.User, error) {
	rows, err := DB.Query(`
        SELECT id, username, email, first_name, last_name, 
               email_verified, user_type, active, created_at
        FROM users 
        WHERE deletion_requested = false
        ORDER BY created_at DESC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(
			&u.ID, &u.Username, &u.Email, &u.FirstName, &u.LastName,
			&u.EmailVerified, &u.UserType, &u.Active, &u.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
