package routes

import (
	"gpsd-api-gateway/internal/gateway/pkg/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", handlers.LogoutHandler).Methods("POST")
	r.HandleFunc("/refresh", handlers.RefreshHandler).Methods("POST")
	// r.HandleFunc("/health", handlers.HealthHandler).Methods("GET")
}
