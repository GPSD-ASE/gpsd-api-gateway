package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gpsd-api-gateway/internal/gateway/pkg/config"

	"github.com/golang-jwt/jwt/v5"
)

var client = &http.Client{
	Timeout: 10 * time.Second,
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type UserData struct {
	Username string `json:"username"`
	Role     string `json:"role,omitempty"`
}

func getUserMgmtBaseURL() string {
	return fmt.Sprintf(
		"http://%s:%s/api/v1",
		config.ApiGatewayConfig.UserMgmtHost,
		config.ApiGatewayConfig.UserMgmtPort,
	)
}

func forwardRequest(w http.ResponseWriter, r *http.Request, endpoint string, modifyBody func([]byte) ([]byte, error)) {
	w.Header().Set("Content-Type", "application/json")

	baseUrl := getUserMgmtBaseURL()

	var err error
	var actualBody []byte

	if r.Body != nil {
		actualBody, err = io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}
		r.Body.Close()
	}

	if modifyBody != nil {
		newBody, err := modifyBody(actualBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}
		actualBody = newBody
	}

	newReq, err := http.NewRequest(r.Method, baseUrl+endpoint, bytes.NewBuffer(actualBody))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	newReq.Header.Set("Content-Type", "application/json")
	if auth := r.Header.Get("Authorization"); auth != "" {
		newReq.Header.Set("Authorization", auth)
	}

	resp, err := client.Do(newReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)

	if _, err := io.Copy(w, resp.Body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	forwardRequest(w, r, "/users", nil)
}

func RegisterAdminHandler(w http.ResponseWriter, r *http.Request) {
	modifyBody := func(original []byte) ([]byte, error) {
		var userData UserData
		if err := json.Unmarshal(original, &userData); err != nil {
			return nil, err
		}
		userData.Role = "admin"
		return json.Marshal(userData)
	}
	forwardRequest(w, r, "/users", modifyBody)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	forwardRequest(w, r, "/signin", nil)
}

// TODO: Remove this from gpsd-api-gateway, only temporary
var secretKey = []byte("secret key")

func VerifyToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return false, err
	}

	if !token.Valid {
		return false, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, fmt.Errorf("invalid token claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return false, fmt.Errorf("invalid expiration claim")
	}

	if time.Now().Unix() > int64(exp) {
		return false, fmt.Errorf("token expired")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return false, fmt.Errorf("invalid username claim")
	}

	return true, nil
}

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	valid := false
	token := r.Header.Get("Authorization")
	if token == "" {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "No token provided"})
		goto out
	}

	valid, err := VerifyToken(token)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		goto out
	}

out:
	w.Header().Set("Content-Type", "application/json")

	if valid {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}

	json.NewEncoder(w).Encode(response)

	// forwardRequest(w, r, "/verify", nil)
}
