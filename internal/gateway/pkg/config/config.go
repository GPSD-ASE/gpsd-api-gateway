package config

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var ApiGatewayConfig *Config

type Config struct {
	VaultAddr        string
	VaultAuth        string
	VaultRole        string
	VaultToken       string
	LogLevel         string
	UserMgmtHost     string
	UserMgmtPort     string
	MapMgmtHost      string
	MapMgmtPort      string
	IncidentMgmtHost string
	IncidentMgmtPort string
}

func findServiceEnvVar(envs []string, service, field string) string {
	suffix := fmt.Sprintf("_%s_%s", service, field)

	for _, env := range envs {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) < 2 {
			continue
		}

		key := pair[0]
		value := pair[1]

		// If this env var ends with  target suffix (e.g., "_USER_MGMT_SERVICE_HOST")
		if strings.HasSuffix(strings.ToUpper(key), suffix) {
			return value
		}
	}

	return ""
}

func LoadConfig(envs []string) {
	ApiGatewayConfig = &Config{
		// Non-service configs remain the same
		VaultAddr:  getEnv("VAULT_ADDR", "http://127.0.0.1:8200"),
		VaultAuth:  getEnv("VAULT_AUTH_METHOD", "token"),
		VaultRole:  getEnv("VAULT_ROLE", ""),
		VaultToken: getEnv("VAULT_TOKEN", ""),
		LogLevel:   getEnv("LOG_LEVEL", "info"),

		// Service configs change dynamically based on release name
		UserMgmtHost:     findServiceEnvVar(envs, "USER_MGMT", "SERVICE_HOST"),
		UserMgmtPort:     findServiceEnvVar(envs, "USER_MGMT", "SERVICE_PORT"),
		MapMgmtHost:      findServiceEnvVar(envs, "MAP_MGMT", "SERVICE_HOST"),
		MapMgmtPort:      findServiceEnvVar(envs, "MAP_MGMT", "SERVICE_PORT"),
		IncidentMgmtHost: findServiceEnvVar(envs, "INCIDENT_MGMT", "SERVICE_HOST"),
		IncidentMgmtPort: findServiceEnvVar(envs, "INCIDENT_MGMT", "SERVICE_PORT"),
	}

	// TODO: Change this to log.Fatal once all deployments are updated.
	if ApiGatewayConfig.UserMgmtHost == "" || ApiGatewayConfig.UserMgmtPort == "" {
		log.Print("User Management service environment variables not found")
	}
	if ApiGatewayConfig.MapMgmtHost == "" || ApiGatewayConfig.MapMgmtPort == "" {
		log.Print("Map Management service environment variables not found")
	}
	if ApiGatewayConfig.IncidentMgmtHost == "" || ApiGatewayConfig.IncidentMgmtPort == "" {
		log.Print("Incident Management service environment variables not found")
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
