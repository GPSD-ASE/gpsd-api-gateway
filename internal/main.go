package main

import (
	"log"
	"net/http"
	"os"

	"gpsd-api-gateway/internal/gateway/pkg/config"
	"gpsd-api-gateway/internal/gateway/pkg/routes"

	"github.com/gorilla/mux"
)

func startHTTPServer() {
	r := mux.NewRouter()
	routes.RegisterRoutes(r)

	server := &http.Server{
		Addr:    ":3000",
		Handler: r,
	}

	log.Println("HTTP  Server for API Gateway running on http://0.0.0.0:3000.")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("HTTP server error:", err)
	}
}

func main() {
	log.Println("Starting API gateway...")

	envs := os.Environ()
	config.LoadConfig(envs)

	startHTTPServer()
}
