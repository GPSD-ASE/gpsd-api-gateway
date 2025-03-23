package routes

import (
	"gpsd-api-gateway/internal/gateway/pkg/config"
	"gpsd-api-gateway/internal/gateway/pkg/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes(cc *config.Config, r *mux.Router) {
	handler := handlers.NewHandler(cc)

	// Health check routes
	r.HandleFunc("/health", handler.NewHealthCheckHandler)
	r.HandleFunc("/ready", handler.NewHealthCheckHandler)

	// User routes
	r.HandleFunc("/register", handler.RegisterHandler).Methods("POST")
	r.HandleFunc("/register-admin", handler.RegisterAdminHandler).Methods("POST")
	r.HandleFunc("/login", handler.LoginHandler).Methods("POST")
	r.HandleFunc("/verify", handler.VerifyHandler).Methods("GET")

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
