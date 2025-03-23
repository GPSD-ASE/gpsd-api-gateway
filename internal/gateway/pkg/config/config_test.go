package config

import (
	"log"
	"os"
	"strings"
	"testing"
)

func TestFindServiceEnvVar(t *testing.T) {
	testCases := []struct {
		name          string
		envs          []string
		service       string
		field         string
		defaultValue  string
		expectedValue string
	}{
		{
			name: "Test 1: Standard Environment Variable",
			envs: []string{
				"GPSD_USER_MGMT_SERVICE_HOST=gpsd-user-service.example.com",
				"GPSD_USER_MGMT_SERVICE_PORT=5500",
			},
			service:       "USER_MGMT",
			field:         "SERVICE_HOST",
			defaultValue:  "localhost",
			expectedValue: "gpsd-user-service.example.com",
		},
		{
			name: "Test 2: Environment Variable With Different Prefix",
			envs: []string{
				"MY_APP_USER_MGMT_SERVICE_HOST=different-prefix.example.com",
				"GPSD_USER_MGMT_SERVICE_PORT=5500",
			},
			service:       "USER_MGMT",
			field:         "SERVICE_HOST",
			defaultValue:  "localhost",
			expectedValue: "different-prefix.example.com",
		},
		{
			name: "Case Insensitive Environment Variable",
			envs: []string{
				"gpsd_user_mgmt_service_host=lowercase.example.com",
				"GPSD_USER_MGMT_SERVICE_PORT=5500",
			},
			service:       "USER_MGMT",
			field:         "SERVICE_HOST",
			defaultValue:  "localhost",
			expectedValue: "lowercase.example.com",
		},
		{
			name: "Missing Environment Variable",
			envs: []string{
				"SOME_OTHER_VAR=value",
			},
			service:       "USER_MGMT",
			field:         "SERVICE_HOST",
			defaultValue:  "default-host",
			expectedValue: "default-host",
		},
		{
			name: "Malformed Environment Variable",
			envs: []string{
				"GPSD_USER_MGMT_SERVICE_HOST",
				"GPSD_USER_MGMT_SERVICE_PORT=5500",
			},
			service:       "USER_MGMT",
			field:         "SERVICE_HOST",
			defaultValue:  "default-for-malformed",
			expectedValue: "default-for-malformed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := FindServiceEnvVar(tc.envs, tc.service, tc.field, tc.defaultValue)
			if result != tc.expectedValue {
				t.Errorf("Expected '%s', got '%s'", tc.expectedValue, result)
			}
		})
	}
}

func TestGetEnvVars(t *testing.T) {
	originalEnv := os.Environ()

	defer func() {
		os.Clearenv()
		for _, env := range originalEnv {
			parts := split2(env, "=")
			if len(parts) == 2 {
				os.Setenv(parts[0], parts[1])
			}
		}
	}()

	testCases := []struct {
		name          string
		key           string
		defaultValue  string
		envVarValue   string
		setEnvVar     bool
		expectedValue string
	}{
		{
			name:          "Environment Variable Exists",
			key:           "TEST_VAR_1",
			defaultValue:  "default1",
			envVarValue:   "value1",
			setEnvVar:     true,
			expectedValue: "value1",
		},
		{
			name:          "Environment Variable Does Not Exist",
			key:           "TEST_VAR_2",
			defaultValue:  "default2",
			setEnvVar:     false,
			expectedValue: "default2",
		},
		{
			name:          "Environment Variable Is Empty",
			key:           "TEST_VAR_3",
			defaultValue:  "default3",
			envVarValue:   "",
			setEnvVar:     true,
			expectedValue: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnvVar {
				os.Setenv(tc.key, tc.envVarValue)
			} else {
				os.Unsetenv(tc.key)
			}

			result := GetEnvVars(tc.key, tc.defaultValue)

			if result != tc.expectedValue {
				t.Errorf("Expected '%s', got '%s'", tc.expectedValue, result)
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	testEnvs := []string{
		"VAULT_ADDR=http://test-vault:8200",
		"VAULT_AUTH_METHOD=kubernetes",
		"VAULT_ROLE=test-role",
		"VAULT_TOKEN=test-token",
		"LOG_LEVEL=debug",
		"API_GATEWAY_APP_PORT=3000",
		"GPSD_USER_MGMT_SERVICE_HOST=user-test-host",
		"GPSD_USER_MGMT_SERVICE_PORT=5501",
		"GPSD_MAP_MGMT_SERVICE_HOST=map-test-host",
		"GPSD_MAP_MGMT_SERVICE_PORT=9001",
		"GPSD_INCIDENT_MGMT_SERVICE_HOST=incident-test-host",
		"GPSD_INCIDENT_MGMT_SERVICE_PORT=7001",
	}

	originalOsEnv := os.Environ()

	os.Clearenv()
	for _, env := range testEnvs {
		parts := split2(env, "=")
		if len(parts) == 2 {
			os.Setenv(parts[0], parts[1])
		}
	}

	originalFatal := logFatal
	fatalCalled := false
	logFatal = func(v ...interface{}) {
		fatalCalled = true
	}

	defer func() {
		logFatal = originalFatal

		os.Clearenv()
		for _, env := range originalOsEnv {
			parts := split2(env, "=")
			if len(parts) == 2 {
				os.Setenv(parts[0], parts[1])
			}
		}
	}()

	cc := LoadConfig(os.Environ())

	if cc.VaultAddr != "http://test-vault:8200" {
		t.Errorf("Expected VaultAddr to be 'http://test-vault:8200', got '%s'", cc.VaultAddr)
	}

	if cc.VaultAuth != "kubernetes" {
		t.Errorf("Expected VaultAuth to be 'kubernetes', got '%s'", cc.VaultAuth)
	}

	if cc.VaultRole != "test-role" {
		t.Errorf("Expected VaultRole to be 'test-role', got '%s'", cc.VaultRole)
	}

	if cc.VaultToken != "test-token" {
		t.Errorf("Expected VaultToken to be 'test-token', got '%s'", cc.VaultToken)
	}

	if cc.LogLevel != "debug" {

		t.Errorf("Expected LogLevel to be 'debug', got '%s'", cc.LogLevel)
	}

	if cc.APIGatewayPort != "3000" {
		t.Errorf("Expected APIGatewayPort to be '3000', got '%s'", cc.APIGatewayPort)
	}

	if cc.UserMgmtHost != "user-test-host" {
		t.Errorf("Expected UserMgmtHost to be 'user-test-host', got '%s'", cc.UserMgmtHost)
	}

	if cc.UserMgmtPort != "5501" {
		t.Errorf("Expected UserMgmtPort to be '5501', got '%s'", cc.UserMgmtPort)
	}

	if cc.MapMgmtHost != "map-test-host" {
		t.Errorf("Expected MapMgmtHost to be 'map-test-host', got '%s'", cc.MapMgmtHost)
	}

	if cc.MapMgmtPort != "9001" {
		t.Errorf("Expected MapMgmtPort to be '9001', got '%s'", cc.MapMgmtPort)
	}

	if cc.IncidentMgmtHost != "incident-test-host" {
		t.Errorf("Expected IncidentMgmtHost to be 'incident-test-host', got '%s'", cc.IncidentMgmtHost)
	}

	if cc.IncidentMgmtPort != "7001" {
		t.Errorf("Expected IncidentMgmtPort to be '7001', got '%s'", cc.IncidentMgmtPort)
	}

	if fatalCalled {
		t.Errorf("Expected log.Fatal not to be called with valid configuration")
	}
}

func TestLoadConfigDefaults(t *testing.T) {
	originalEnv := os.Environ()

	defer func() {
		os.Clearenv()
		for _, env := range originalEnv {
			parts := split2(env, "=")
			if len(parts) == 2 {
				os.Setenv(parts[0], parts[1])
			}
		}
	}()

	cc := LoadConfig(os.Environ())

	if cc.VaultAddr != "http://127.0.0.1:8200" {
		t.Errorf("Expected default VaultAddr to be 'http://127.0.0.1:8200', got '%s'", cc.VaultAddr)
	}

	if cc.VaultAuth != "token" {
		t.Errorf("Expected default VaultAuth to be 'token', got '%s'", cc.VaultAuth)
	}

	if cc.VaultRole != "kubernetes" {
		t.Errorf("Expected default VaultRole to be 'kubernetes', got '%s'", cc.VaultRole)
	}

	if cc.VaultToken != "root" {
		t.Errorf("Expected default VaultToken to be 'root', got '%s'", cc.VaultToken)
	}

	if cc.LogLevel != "info" {
		t.Errorf("Expected default LogLevel to be 'info', got '%s'", cc.LogLevel)
	}

	if cc.APIGatewayPort != "3000" {
		t.Errorf("Expected default APIGatewayPort to be '3000', got '%s'", cc.APIGatewayPort)
	}

	if cc.UserMgmtHost != "localhost" {
		t.Errorf("Expected default UserMgmtHost to be 'localhost', got '%s'", cc.UserMgmtHost)
	}

	if cc.UserMgmtPort != "5500" {
		t.Errorf("Expected default UserMgmtPort to be '5500', got '%s'", cc.UserMgmtPort)
	}

	if cc.MapMgmtHost != "localhost" {
		t.Errorf("Expected default MapMgmtHost to be 'localhost', got '%s'", cc.MapMgmtHost)
	}

	if cc.MapMgmtPort != "9000" {
		t.Errorf("Expected default MapMgmtPort to be '9000', got '%s'", cc.MapMgmtPort)
	}

	if cc.IncidentMgmtHost != "localhost" {
		t.Errorf("Expected default IncidentMgmtHost to be 'localhost', got '%s'", cc.IncidentMgmtHost)
	}

	if cc.IncidentMgmtPort != "7000" {
		t.Errorf("Expected default IncidentMgmtPort to be '7000', got '%s'", cc.IncidentMgmtPort)
	}
}

func split2(s, sep string) []string {
	parts := make([]string, 2)
	i := 0
	if i = len(s); i > 0 {
		if i = strings.Index(s, sep); i >= 0 {
			parts[0] = s[:i]
			if i+len(sep) < len(s) {
				parts[1] = s[i+len(sep):]
			}
			return parts
		}

	}
	parts[0] = s
	return parts[:1]
}

var logFatal = log.Fatal
