package middlewares

import (
	"context"
	"mental-health/config"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

const RoleKey ContextKey = "role"

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: no token provided", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Unauthorized: invalid token format", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return config.JWTSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized: invalid token claims", http.StatusUnauthorized)
			return
		}

		role, ok := claims["role"].(string)
		if !ok || role != "admin" {
			http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "Unauthorized: invalid user ID", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, RoleKey, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
