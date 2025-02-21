package routes

import (
	"gpsd-api-gateway/internal/gateway/pkg/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST") // forward to user service
	r.HandleFunc("/register-admin", handlers.RegisterAdminHandler).Methods("POST") // forward to user service, add roleType in body
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST") // -> forward to user service, return the token from user service
	r.HandleFunc("/verify", handlers.VerifyHandler).Methods("POST") // -> for all users, logic in api-gateway
}
