package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"gpsd-api-gateway/internal/gateway/pkg/handlers"
	"gpsd-api-gateway/internal/gateway/pkg/routes"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("Starting API gateway...")
	cert, err := handlers.RetrieveCertFromVault()
	if err != nil {
		log.Fatalf("Error retrieving certificates from Vault: %v", err)
	}

	server := &http.Server{
		Addr:    ":3000",
		Handler: mux.NewRouter(),
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
	}

	r := mux.NewRouter()
	routes.RegisterRoutes(r)

	certFile, keyFile, err := handlers.WriteCertificateAndKey(*cert)
	if err != nil {
		log.Fatalf("Error creating temp key cert file: %w", err)
	}

	// Start the server with the certificate loaded from Vault.
	log.Println("API Gateway running on https://localhost:3000")
	log.Fatal(server.ListenAndServeTLS(certFile, keyFile))
}
