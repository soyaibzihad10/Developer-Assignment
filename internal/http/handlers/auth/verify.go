package auth

import (
    "net/http"

    "github.com/soyaibzihad10/Developer-Assignment/internal/database"
)

func VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
    token := r.URL.Query().Get("token")
    if token == "" {
        http.Error(w, "Token is required", http.StatusBadRequest)
        return
    }

    query := "UPDATE users SET email_verified = true WHERE verification_token = $1 AND token_expiry > NOW()"
    result, err := database.DB.Exec(query, token)
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil || rowsAffected == 0 {
        http.Error(w, "Invalid or expired token", http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Email verified successfully"))
}