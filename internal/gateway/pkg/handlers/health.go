package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gpsd-api-gateway/internal/gateway/pkg/config"
)

type VaultHealth struct {
	Initialized bool   `json:"initialized"`
	Sealed      bool   `json:"sealed"`
	Standby     bool   `json:"standby"`
	Error       string `json:"error,omitempty"`
}
type HealthResponse struct {
	Status      string       `json:"status"`
	Timestamp   string       `json:"timestamp"`
	VaultStatus *VaultHealth `json:"vault_status,omitempty"`
}

func checkVaultHealth(vaultAddr string) *VaultHealth {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(fmt.Sprintf("%s/v1/sys/health", vaultAddr))
	if err != nil {
		return &VaultHealth{
			Error: fmt.Sprintf("Failed to connect to Vault: %v", err),
		}
	}
	defer resp.Body.Close()

	var health VaultHealth
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		health.Error = fmt.Sprintf("Failed to decode Vault response: %v", err)
	}

	return &health
}

// HealthCheckHandler for HTTP-only endpoints.
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	cfg := config.LoadConfig()

	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	vaultHealth := checkVaultHealth(cfg.VaultAddr)
	if vaultHealth.Error != "" || vaultHealth.Sealed {
		response.Status = "degraded"
	}
	response.VaultStatus = vaultHealth

	if response.Status == "ok" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
