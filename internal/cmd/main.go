package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/GPSD-ASE/gpsd-api-gateway//pkg/routes"
	"github.com/GPSD-ASE/gpsd-api-gateway/pkg/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// Load the certificate and key from Vault.
	cert, err := handlers.RetrieveCertFromVault()
	if err != nil {
		log.Fatalf("Error retrieving certificates from Vault: %v", err)
	}

	// Create the HTTPS server using the retrieved certificate.
	server := &http.Server{
		Addr:    ":3000",
		Handler: mux.NewRouter(),
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
	}

	// Register routes
	r := mux.NewRouter()
	routes.RegisterRoutes(r)

	// Start the server with the certificate loaded from Vault
	log.Println("API Gateway running on https://localhost:3000")
	log.Fatal(server.ListenAndServeTLS(cert.Certificate[0], cert.PrivateKey))

}
