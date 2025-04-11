package routes

import (
	"gpsd-api-gateway/internal/gateway/pkg/config"
	"gpsd-api-gateway/internal/gateway/pkg/handlers"
	"gpsd-api-gateway/internal/gateway/pkg/middleware"

	"github.com/gorilla/mux"
)

func RegisterRoutes(cc *config.Config, r *mux.Router) {

	r.Use(middleware.RequestLogger)
	// r.Use(middleware.RequestLogger, middleware.AuthMiddleware)

	handler := handlers.NewHandler(cc)

	// Health check routes
	r.HandleFunc("/health", handler.NewHealthCheckHandler)
	r.HandleFunc("/ready", handler.NewHealthCheckHandler)

	// Auth routes
	r.HandleFunc("/register", handler.RegisterHandler).Methods("POST")
	r.HandleFunc("/register-admin", handler.RegisterAdminHandler).Methods("POST")
	r.HandleFunc("/signin", handler.SigninHandler).Methods("POST")
	r.HandleFunc("/signout", handler.SignoutHandler).Methods("POST")
	r.HandleFunc("/verify", handler.VerifyHandler).Methods("GET")

	// User routes
	r.HandleFunc("/users", handler.GetUsersHandler).Methods("GET")
	r.HandleFunc("/users/{id}", handler.GetUserByIdHandler).Methods("GET")
	r.HandleFunc("/users/{id}", handler.UpdateUserHandler).Methods("PATCH")
	r.HandleFunc("/users/{id}", handler.DeleteUserHandler).Methods("DELETE")

	// Map routes
	r.HandleFunc("/zones", handler.ZonesHandler).Methods("GET")
	r.HandleFunc("/routing", handler.RoutingHandler).Methods("GET")
	r.HandleFunc("/evacuation", handler.EvacuationHandler).Methods("POST")
	r.HandleFunc("/traffic", handler.TrafficHandler).Methods("GET")

	// Incident routes
	r.HandleFunc("/incidents", handler.GetAllIncidentsHandler).Methods("GET")
	r.HandleFunc("/incidents", handler.CreateIncidentHandler).Methods("POST")
	r.HandleFunc("/incidents/{id}", handler.GetIncidentByIdHandler).Methods("GET")
	r.HandleFunc("/incidents/{id}", handler.DeleteIncidentHandler).Methods("DELETE")
	r.HandleFunc("/incidents/{id}/status/{status}", handler.ChangeIncidentStatusHandler).Methods("PATCH")
}
