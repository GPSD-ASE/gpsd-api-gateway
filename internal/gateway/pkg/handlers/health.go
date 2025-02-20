package handlers

import (
	"fmt"
	"net/http"
)

// HealthCheckHandler for HTTP-only endpoints.
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Health: OK")
}