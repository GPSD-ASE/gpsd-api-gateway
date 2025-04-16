package config

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Config struct {
	VaultAddr          string
	VaultAuth          string
	VaultRole          string
	VaultToken         string
	LogLevel           string
	APIGatewayPort     string
	UserMgmtHost       string
	UserMgmtPort       string
	MapMgmtHost        string
	MapMgmtPort        string
	IncidentMgmtHost   string
	IncidentMgmtPort   string
	DecisionEngineHost string
	DecisionEnginePort string
	EscalationMgmtHost string
	EscalationMgmtPort string
}

func FindServiceEnvVar(envs []string, service, field string, def string) string {
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

	return def
}

func LoadConfig(envs []string) *Config {
	cc := &Config{
		// Non-service configs remain the same
		VaultAddr:  GetEnvVars("VAULT_ADDR", "http://127.0.0.1:8200"),
		VaultAuth:  GetEnvVars("VAULT_AUTH_METHOD", "token"),
		VaultRole:  GetEnvVars("VAULT_ROLE", "kubernetes"),
		VaultToken: GetEnvVars("VAULT_TOKEN", "root"),
		LogLevel:   GetEnvVars("LOG_LEVEL", "info"),

		APIGatewayPort: GetEnvVars("API_GATEWAY_APP_PORT", "3000"),

		// Service configs change dynamically based on release name
		UserMgmtHost:       FindServiceEnvVar(envs, "USER_MGMT", "SERVICE_HOST", "localhost"),
		UserMgmtPort:       FindServiceEnvVar(envs, "USER_MGMT", "SERVICE_PORT", "5500"),
		MapMgmtHost:        FindServiceEnvVar(envs, "MAP_MGMT", "SERVICE_HOST", "localhost"),
		MapMgmtPort:        FindServiceEnvVar(envs, "MAP_MGMT", "SERVICE_PORT", "9000"),
		IncidentMgmtHost:   FindServiceEnvVar(envs, "INCIDENT_MGMT", "SERVICE_HOST", "localhost"),
		IncidentMgmtPort:   FindServiceEnvVar(envs, "INCIDENT_MGMT", "SERVICE_PORT", "7000"),
		DecisionEngineHost: FindServiceEnvVar(envs, "RESOURCE_MAPPING", "SERVICE_HOST", "localhost"),
		DecisionEnginePort: FindServiceEnvVar(envs, "RESOURCE_MAPPING", "SERVICE_PORT", "7134"),
		EscalationMgmtHost: FindServiceEnvVar(envs, "ESCALATION_MGMT", "SERVICE_HOST", "localhost"),
		EscalationMgmtPort: FindServiceEnvVar(envs, "ESCALATION_MGMT", "SERVICE_PORT", "8000"),
	}

	if cc.UserMgmtHost == "" || cc.UserMgmtPort == "" {
		log.Fatal("User Management service environment variables not found")
	}
	if cc.MapMgmtHost == "" || cc.MapMgmtPort == "" {
		log.Fatal("Map Management service environment variables not found")
	}
	if cc.IncidentMgmtHost == "" || cc.IncidentMgmtPort == "" {
		log.Fatal("Incident Management service environment variables not found")
	}

	return cc
}

func GetEnvVars(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
