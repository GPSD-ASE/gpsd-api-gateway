package config

import (
	"fmt"
	"log"
	"os"
	"strings"
)

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

func LoadConfig() *Config {
	envs := os.Environ()

	cfg := &Config{
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

	if cfg.UserMgmtHost == "" || cfg.UserMgmtPort == "" {
		log.Fatal("User Management service environment variables not found")
	}
	if cfg.MapMgmtHost == "" || cfg.MapMgmtPort == "" {
		log.Fatal("Map Management service environment variables not found")
	}
	if cfg.IncidentMgmtHost == "" || cfg.IncidentMgmtPort == "" {
		log.Fatal("Incident Management service environment variables not found")
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
