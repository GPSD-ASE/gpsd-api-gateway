package routes

import (
	"gpsd-api-gateway/internal/gateway/pkg/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	// Health check routes
	r.HandleFunc("/health", handlers.HealthCheckHandler)
	r.HandleFunc("/ready", handlers.HealthCheckHandler)

	// User routes
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/register-admin", handlers.RegisterAdminHandler).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/verify", handlers.VerifyHandler).Methods("GET")

	// Map routes
	r.HandleFunc("/zones", handlers.ZonesHandler).Methods("GET")
	r.HandleFunc("/routing", handlers.RoutingHandler).Methods("GET")
	r.HandleFunc("/evacuation", handlers.EvacuationHandler).Methods("POST")
	r.HandleFunc("/traffic", handlers.TrafficHandler).Methods("GET")

	// Incident routes
	r.HandleFunc("/incidents", handlers.GetAllIncidentsHandler).Methods("GET")
	r.HandleFunc("/incidents", handlers.CreateIncidentHandler).Methods("POST")
	r.HandleFunc("/incidents/{id}", handlers.GetIncidentByIdHandler).Methods("GET")
	r.HandleFunc("/incidents/{id}", handlers.DeleteIncidentHandler).Methods("DELETE")
	r.HandleFunc("/incidents/{id}/status/{status}", handlers.ChangeIncidentStatusHandler).Methods("PATCH")
}
