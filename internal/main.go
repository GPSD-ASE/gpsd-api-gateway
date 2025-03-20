package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"gpsd-api-gateway/internal/gateway/pkg/config"
	"gpsd-api-gateway/internal/gateway/pkg/handlers"
	"gpsd-api-gateway/internal/gateway/pkg/routes"

	"github.com/gorilla/mux"
)

// Start HTTPS Server for handling external requests.
func startHTTPSServer() {
	cert, err := handlers.RetrieveCertFromVault()
	if err != nil {
		log.Fatalf("Error retrieving certificates from Vault: %v", err)
	}

	r := mux.NewRouter()
	routes.RegisterRoutes(r)

	server := &http.Server{
		Addr:    ":3000",
		Handler: r,
		TLSConfig: &tls.Config{
			MinVersion:   tls.VersionTLS13,
			Certificates: []tls.Certificate{*cert},
		},
	}

	log.Println("HTTPS Server for API Gateway running on https://0.0.0.0:3000.")
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatalf("HTTPS server error: %v", err)
	}
}

// Start HTTP Server
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

	log.Printf("Vault is running at %s.", config.ApiGatewayConfig.VaultAddr)

	startHTTPServer()
	// startHTTPSServer()
}
