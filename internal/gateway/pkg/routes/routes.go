package routes

import (
	"gpsd-api-gateway/internal/gateway/pkg/config"
	"gpsd-api-gateway/internal/gateway/pkg/handlers"
	"gpsd-api-gateway/internal/gateway/pkg/middleware"

	"github.com/gorilla/mux"
)

func RegisterRoutes(cc *config.Config, r *mux.Router) {

	r.Use(middleware.RequestLogger)

	handler := handlers.NewHandler(cc)

	/* Public routes do not require authentication */
	publicRoutes := r.PathPrefix("").Subrouter()

	/* Health check routes */
	publicRoutes.HandleFunc("/health", handler.NewHealthCheckHandler)
	publicRoutes.HandleFunc("/ready", handler.NewHealthCheckHandler)

	/* Auth routes */
	publicRoutes.HandleFunc("/register", handler.RegisterHandler).Methods("POST")
	publicRoutes.HandleFunc("/register-admin", handler.RegisterAdminHandler).Methods("POST")
	publicRoutes.HandleFunc("/signin", handler.SigninHandler).Methods("POST")
	publicRoutes.HandleFunc("/verify", handler.VerifyHandler).Methods("GET")

	/* Protected routes require authentication */
	protectedRoutes := r.PathPrefix("").Subrouter()
	protectedRoutes.Use(middleware.AuthMiddleware)

	protectedRoutes.HandleFunc("/signout", handler.SignoutHandler).Methods("POST")

	/* User routes */
	protectedRoutes.HandleFunc("/users", handler.GetUsersHandler).Methods("GET")
	protectedRoutes.HandleFunc("/users/{id}", handler.GetUserByIdHandler).Methods("GET")
	protectedRoutes.HandleFunc("/users/{id}", handler.UpdateUserHandler).Methods("PATCH")
	protectedRoutes.HandleFunc("/users/{id}", handler.DeleteUserHandler).Methods("DELETE")

	/* Map routes */
	protectedRoutes.HandleFunc("/zones", handler.GetZonesHandler).Methods("GET")
	protectedRoutes.HandleFunc("/route", handler.GetRouteHandler).Methods("GET")
	protectedRoutes.HandleFunc("/routing", handler.GetRoutingHandler).Methods("GET")
	protectedRoutes.HandleFunc("/safezones", handler.GetSafezonesHandler).Methods("GET")
	protectedRoutes.HandleFunc("/safezones", handler.PostSafezonesHandler).Methods("POST")
	protectedRoutes.HandleFunc("/evacuation", handler.PostEvacuationHandler).Methods("POST")
	protectedRoutes.HandleFunc("/traffic", handler.TrafficHandler).Methods("GET")

	/* Incident routes */
	protectedRoutes.HandleFunc("/incidents", handler.GetAllIncidentsHandler).Methods("GET")
	protectedRoutes.HandleFunc("/incidents", handler.PostIncidentHandler).Methods("POST")
	protectedRoutes.HandleFunc("/incidents/{id}", handler.GetIncidentByIdHandler).Methods("GET")
	protectedRoutes.HandleFunc("/incidents/{id}", handler.DeleteIncidentHandler).Methods("DELETE")
	protectedRoutes.HandleFunc("/incidents/{id}/status/{status}", handler.UpdateIncidentStatusHandler).Methods("PATCH")

	/* Decision and escalation engine routes */
	protectedRoutes.HandleFunc("/decision/incident", handler.PostDecisionHandler).Methods("POST")
	protectedRoutes.HandleFunc("/incident-analysis", handler.PostIncidentAnalysis).Methods("GET")
}
