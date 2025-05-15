package middleware

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// check if the token is valid
		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// get the claims from the token
		claims := token.Claims.(jwt.MapClaims)
		r.Header.Set("user_id", claims["user_id"].(string))
		r.Header.Set("user_type", claims["user_type"].(string))

		next.ServeHTTP(w, r)
	})
}
