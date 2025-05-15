package middleware

import (
    "net/http"
    "os"
    "github.com/golang-jwt/jwt/v4"
)

// IsAuthenticated checks if the user is authenticated by verifying the JWT token in the cookie
func IsAuthenticated(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // take the JWT token from the cookie
        cookie, err := r.Cookie("auth_token")
        if err != nil {
            http.Error(w, "invalid access", http.StatusUnauthorized)
            return
        }

        // check if the token is valid
        token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("JWT_SECRET")), nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "invalid token", http.StatusUnauthorized)
            return
        }

        // take the claims from the token
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            http.Error(w, "invalid claims", http.StatusUnauthorized)
            return
        }

        // add user information to the request header
        r.Header.Set("user_id", claims["user_id"].(string))
        r.Header.Set("user_email", claims["email"].(string))

        next.ServeHTTP(w, r)
    })
}