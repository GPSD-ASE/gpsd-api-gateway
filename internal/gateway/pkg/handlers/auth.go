package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// TODO: Perform authentication (JWT creation or other logic here).
	if creds.Username == "admin" && creds.Password == "password" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Login successful")
		return
	}

	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Invalidate the JWT token.
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Logged out successfully")
}

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Generate a new JWT token based on the refresh token.
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Token refreshed")
}
