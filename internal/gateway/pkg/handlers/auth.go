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
	var response Response
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		response = Response{Message: "Invalid JSON payload from user."}
		w.WriteHeader(http.StatusBadRequest)
		goto out
	}

	if loginRequest.Username == "user1" {
		response = Response{Message: "Login forbidden!"}
		w.WriteHeader(http.StatusForbidden)
		goto out
	}

	response = Response{Message: "Login API is successful.", AccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjEwIiwidXNlcm5hbWUiOiJEdW1teSIsInJvbGUiOiIxIn0.isgyco7Uq5J4N-QIFW3RJ_JM7eD4WMAJoqDh0bxVrYo "}
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
