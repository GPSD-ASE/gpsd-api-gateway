package middleware

import (
	"encoding/json"
	"gpsd-api-gateway/internal/gateway/pkg/handlers"
	"log"
	"net/http"
	"strings"
	"time"
)

/* RequestLogger logs information about incoming HTTP requests. */
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		log.Printf("REQUEST: %s %s %s | Remote: %s | User-Agent: %s",
			r.Method,
			r.URL.Path,
			r.Proto,
			r.RemoteAddr,
			r.Header.Get("User-Agent"),
		)

		next.ServeHTTP(w, r)

		log.Printf("RESPONSE: %s %s | Duration: %v",
			r.Method,
			r.URL.Path,
			time.Since(startTime),
		)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(handlers.ErrorResponse{Error: "no token provided"})
			return
		}

		token := ""
		if len(authHeader) > 7 && strings.HasPrefix(authHeader, "Bearer ") {
			token = authHeader[7:]
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(handlers.ErrorResponse{Error: "invalid authorization format"})
			return
		}

		valid, err := handlers.VerifyToken(token)
		if !valid || err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(handlers.ErrorResponse{Error: err.Error()})
			return
		}

		next.ServeHTTP(w, r)
	})
}
