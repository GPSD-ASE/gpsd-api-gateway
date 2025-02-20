package routes

import (
	"gpsd-api-gateway/internal/gateway/pkg/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST") // non admin users (mobile) -> user service
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST") // -> user service
	r.HandleFunc("/logout", handlers.LogoutHandler).Methods("POST") // -> user service
}
