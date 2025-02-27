package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gpsd-api-gateway/internal/gateway/pkg/config"

	"github.com/golang-jwt/jwt/v5"
)

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

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	ForwardRequest(w, r, getUserMgmtBaseURL()+"/users", nil)
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
	ForwardRequest(w, r, getUserMgmtBaseURL()+"/users", modifyBody)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	ForwardRequest(w, r, getUserMgmtBaseURL()+"/signin", nil)
}

// TODO: Remove this from gpsd-api-gateway, only temporary
var secretKey = []byte("a-string-secret-at-least-256-bits-long")

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

	return true, nil
}

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "no token provided"})
		return
	}

	token := ""
	if len(authHeader) > 7 && strings.HasPrefix(authHeader, "Bearer ") {
		token = authHeader[7:]
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid authorization format"})
		return
	}

	valid, err := VerifyToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]bool{"valid": false})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"valid": true})
}
