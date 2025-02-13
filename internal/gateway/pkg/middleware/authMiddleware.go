package middleware

import (
	"net/http"
	"strings"
)

// Simple authentication middleware example to validate the presence of a token
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// TODO: can validate the JWT token

		next.ServeHTTP(w, r)
	})
}
