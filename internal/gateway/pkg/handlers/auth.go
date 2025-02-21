package handlers

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	AccessToken string `json:accessToken,omitempty`
}

type LoginRequest struct {
	Username string `json:"username"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := Response{Message: "Register now as User! Work in progress."}
	json.NewEncoder(w).Encode(response)
}

func RegisterAdminHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := Response{Message: "Register Now as Admin! Work in progress."}
	json.NewEncoder(w).Encode(response)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var loginRequest LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		response := Response{Message: "Invalid JSON payload from user."}
		w.WriteHeader(http.StatusBadRequest)
		goto out
	}

	if loginRequest.Username == "user1" {
		http.Error(w, "Login forbidden!", http.StatusForbidden)
		response := Response{Message: "Login forbidden!"}
		w.WriteHeader(http.StatusForbidden)
		goto out
	}

	response := Response{Message: "Login is successful."}
	w.WriteHeader(http.StatusOK)

out:
	json.NewEncoder(w).Encode(response)
}

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := Response{Message: "Token verification is successful."}
	json.NewEncoder(w).Encode(response)
}
