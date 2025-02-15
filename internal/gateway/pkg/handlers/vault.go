package handlers

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"net/http"

	"github.com/hashicorp/vault/api"
)

/*
WriteCertificateAndKey writes the certificate and key to temporary files
and returns the file paths for further use. The temporary files will be deleted
when the function exits.
*/
func WriteCertificateAndKey(cert tls.Certificate) (string, string, error) {
	certFile, err := os.CreateTemp("", "cert.pem")
	if err != nil {
		return "", "", fmt.Errorf("error creating temp cert file: %w", err)
	}
	defer func() {
		certFile.Close()
		os.Remove(certFile.Name())
	}()

	if _, err := certFile.Write(cert.Certificate[0]); err != nil {
		return "", "", fmt.Errorf("error writing to temp cert file: %w", err)
	}

	keyFile, err := os.CreateTemp("", "key.pem")
	if err != nil {
		return "", "", fmt.Errorf("error creating temp key file: %w", err)
	}
	defer func() {
		keyFile.Close()
		os.Remove(keyFile.Name())
	}()

	keyBytes, ok := cert.PrivateKey.([]byte)
	if !ok {
		return "", "", fmt.Errorf("error asserting type of private key")
	}

	if _, err := keyFile.Write(keyBytes); err != nil {
		return "", "", fmt.Errorf("error writing to temp key file: %w", err)
	}

	return certFile.Name(), keyFile.Name(), nil
}

/*	RetrieveCertFromVault retrieves certificate and key from Vault. */
func RetrieveCertFromVault() (*tls.Certificate, error) {
	// Load the vault CA certificate.
	caCert, err := os.ReadFile("/etc/ssl/certs/vault.pem")
	if err != nil {
		return nil, fmt.Errorf("unable to read custom CA certificate: %v", err)
	}

	// Create a certificate pool and add the vault CA.
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	
	vaultAddr := os.Getenv("VAULT_ADDR")
	if vaultAddr == "" {
		return nil, fmt.Errorf("VAULT_ADDR is not set.")
	}

	client, err := api.NewClient(&api.Config{
		Address: vaultAddr,
		HttpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: caCertPool,
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create Vault client: %v", err)
	}

	vaultToken := os.Getenv("VAULT_TOKEN")
	if vaultToken == "" {
		return nil, fmt.Errorf("VAULT_TOKEN is not set.")
	}
	client.SetToken(vaultToken)

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
