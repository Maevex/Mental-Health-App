package middlewares

import (
	"context"
	"mental-health/config"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// Context key
type ContextKey string

const UserIDKey ContextKey = "user_id"

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ambil token dari header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
			return
		}

		// Format token: "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Unauthorized: Invalid token format", http.StatusUnauthorized)
			return
		}

		// Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return config.JWTSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// Ambil user_id dari token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized: Invalid token claims", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user_id"].(float64) // JWT mengembalikan angka sebagai float64
		if !ok {
			http.Error(w, "Unauthorized: Invalid user_id", http.StatusUnauthorized)
			return
		}

		// Simpan user_id ke context
		ctx := context.WithValue(r.Context(), UserIDKey, int(userID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
