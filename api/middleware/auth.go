package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"

	"department-management/utils"
)

// AuthMiddleware checks for a valid JWT token and verifies the user role.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if JWT token was provided
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("Authorization header missing")
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Extract token string
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		log.Printf("Extracted token: %s", tokenString)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(utils.GetJWTSecret()), nil
		})

		// Validate token
		if err != nil || !token.Valid {
			log.Printf("Invalid token: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("Failed to parse token claims")
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Check for user role
		role, ok := claims["role"].(string)
		if !ok {
			log.Println("Role claim missing or invalid")
			http.Error(w, "Unauthorized. Role claim missing or invalid", http.StatusUnauthorized)
			return
		}

		// Routes can only be accessed by admin
		if role != "admin" {
			log.Println("User is not an admin")
			http.Error(w, "Forbidden. User is not an admin", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
