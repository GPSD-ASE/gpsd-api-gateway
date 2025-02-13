package handlers

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
)

// Retrieve certificate and key from Vault.
func RetrieveCertFromVault() (*tls.Certificate, error) {
	client, err := api.NewClient(&api.Config{
		Address: "http://localhost:8200",
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create Vault client: %v", err)
	}

	client.SetToken(os.Getenv("VAULT_TOKEN"))

	// Fetch the certificate and key from Vault.
	secret, err := client.KVv2("secret").Get(context.Background(), "api-gateway/cert")
	if err != nil {
		return nil, fmt.Errorf("unable to read secret from Vault: %v", err)
	}

	// Retrieve the certificate and key data.
	certPEM := []byte(secret.Data["cert"].(string))
	keyPEM := []byte(secret.Data["key"].(string))

	// Create the TLS certificate from PEM data.
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate and key: %v", err)
	}

	return &cert, nil
}
