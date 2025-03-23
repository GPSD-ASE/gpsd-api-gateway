package handlers

import (
	"reflect"
	"testing"

	"gpsd-api-gateway/internal/gateway/pkg/config"
)

/*
 * TestNewHandler tests that the NewHandler function
 * correctly initializes a Handler with the provided config.
 */
func TestNewHandler(t *testing.T) {
	testConfig := &config.Config{
		VaultAddr:        "http://localhost:8200",
		VaultAuth:        "test-vault-auth",
		VaultRole:        "test-vault-role",
		VaultToken:       "test-vault-token",
		LogLevel:         "info",
		APIGatewayPort:   "3000",
		UserMgmtHost:     "gpsd-user-service",
		UserMgmtPort:     "8080",
		MapMgmtHost:      "gpsd-map-service",
		MapMgmtPort:      "8081",
		IncidentMgmtHost: "gpsd-incident-service",
		IncidentMgmtPort: "8082",
	}

	handler := NewHandler(testConfig)

	if handler == nil {
		t.Fatal("Expected handler to be created, got nil")
	}

	if handler.Config == nil {
		t.Fatal("Expected handler.Config to be set, got nil")
	}

	if !reflect.DeepEqual(handler.Config, testConfig) {
		t.Errorf("Handler config does not match. Expected %+v, got %+v", testConfig, handler.Config)
	}
}

/*
 * TestNewHandlerWithNilConfig tests that
 * NewHandler handles a nil config appropriately.
 */
func TestNewHandlerWithNilConfig(t *testing.T) {
	handler := NewHandler(nil)

	if handler == nil {
		t.Fatal("Expected handler to be created even with nil config, got nil")
	}

	if handler.Config != nil {
		t.Errorf("Expected handler.Config to be nil, got %+v", handler.Config)
	}
}

/*
 * TestHandlerWithDefaultConfig tests
 * creating a handler with a default config.
 */
func TestHandlerWithDefaultConfig(t *testing.T) {
	defaultConfig := &config.Config{
		VaultAddr:      "http://vault:8200",
		VaultAuth:      "default-vault-auth",
		VaultToken:     "default-vault-token",
		LogLevel:       "info",
		APIGatewayPort: "3030",
	}

	handler := NewHandler(defaultConfig)

	if !reflect.DeepEqual(handler.Config, defaultConfig) {
		t.Errorf("Handler config does not match. Expected %+v, got %+v", defaultConfig, handler.Config)
	}
}

/*
 * TestHandlerConfigAccess tests
 * accessing the config from the handler.
 */
func TestHandlerConfigAccess(t *testing.T) {
	testConfig := &config.Config{
		VaultAddr:        "http://localhost:8200",
		VaultAuth:        "test-vault-auth-2",
		VaultRole:        "test-vault-role-2",
		VaultToken:       "test-vault-token-2",
		LogLevel:         "debug",
		APIGatewayPort:   "3000",
		UserMgmtHost:     "gpsd-user-service",
		UserMgmtPort:     "8080",
		MapMgmtHost:      "gpsd-map-service",
		MapMgmtPort:      "8081",
		IncidentMgmtHost: "gpsd-incident-service",
		IncidentMgmtPort: "8082",
	}

	handler := NewHandler(testConfig)

	if handler.Config.VaultAddr != testConfig.VaultAddr {
		t.Errorf("Expected VaultAddr to be %s, got %s", testConfig.VaultAddr, handler.Config.VaultAddr)
	}

	if handler.Config.LogLevel != testConfig.LogLevel {
		t.Errorf("Expected LogLevel to be %s, got %s", testConfig.LogLevel, handler.Config.LogLevel)
	}

	if handler.Config.APIGatewayPort != testConfig.APIGatewayPort {
		t.Errorf("Expected APIGatewayPort to be %s, got %s", testConfig.APIGatewayPort, handler.Config.APIGatewayPort)
	}
}
