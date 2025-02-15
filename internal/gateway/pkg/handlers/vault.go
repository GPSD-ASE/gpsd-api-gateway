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
		return "", "", fmt.Errorf("Error creating temp cert file: %w", err)
	}
	defer func() {
		certFile.Close()
		os.Remove(certFile.Name())
	}()

	if _, err := certFile.Write(cert.Certificate[0]); err != nil {
		return "", "", fmt.Errorf("Error writing to temp cert file: %w", err)
	}

	keyFile, err := os.CreateTemp("", "key.pem")
	if err != nil {
		return "", "", fmt.Errorf("Error creating temp key file: %w", err)
	}
	defer func() {
		keyFile.Close()
		os.Remove(keyFile.Name())
	}()

	var keyPEM []byte
	switch key := cert.PrivateKey.(type) {
	case *rsa.PrivateKey:
		keyPEM = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	default:
		return "", "", fmt.Errorf("Unsupported private key type: %T", cert.PrivateKey)
	}

	if _, err := keyFile.Write(keyPEM); err != nil {
		return "", "", fmt.Errorf("Error writing to temp key file: %w", err)
	}

	if _, err := keyFile.Write(keyBytes); err != nil {
		return "", "", fmt.Errorf("Error writing to temp key file: %w", err)
	}

	return certFile.Name(), keyFile.Name(), nil
}

/*	RetrieveCertFromVault retrieves certificate and key from Vault. */
func RetrieveCertFromVault() (*tls.Certificate, error) {
	// Load the vault CA certificate.
	caCert, err := os.ReadFile("/etc/ssl/certs/vault.pem")
	if err != nil {
		return nil, fmt.Errorf("Unable to read custom CA certificate: %v", err)
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
		return nil, fmt.Errorf("Unable to create Vault client: %v", err)
	}

	vaultToken := os.Getenv("VAULT_TOKEN")
	if vaultToken == "" {
		return nil, fmt.Errorf("VAULT_TOKEN is not set.")
	}
	client.SetToken(vaultToken)

	// Fetch the certificate and key from Vault.
	secret, err := client.KVv2("secret").Get(context.Background(), "api-gateway/cert")
	if err != nil {
		return nil, fmt.Errorf("Unable to read secret from Vault: %v", err)
	}

	// Retrieve the certificate and key data.
	certPEM := []byte(secret.Data["cert"].(string))
	keyPEM := []byte(secret.Data["key"].(string))

	// Create the TLS certificate from PEM data.
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse certificate and key: %v", err)
	}

	return &cert, nil
}
