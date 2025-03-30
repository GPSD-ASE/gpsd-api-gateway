package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gpsd-api-gateway/internal/gateway/pkg/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type UserData struct {
	Username string `json:"username"`
	Role     string `json:"role,omitempty"`
}

func getUserMgmtBaseURL(cc *config.Config) string {
	return fmt.Sprintf(
		"http://%s:%s/api/v1",
		cc.UserMgmtHost,
		cc.UserMgmtPort,
	)
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	ForwardRequest(w, r, getUserMgmtBaseURL(h.Config)+"/users", nil)
}

func (h *Handler) RegisterAdminHandler(w http.ResponseWriter, r *http.Request) {
	modifyBody := func(original []byte) ([]byte, error) {
		var userData UserData
		if err := json.Unmarshal(original, &userData); err != nil {
			return nil, err
		}
		userData.Role = "2"
		return json.Marshal(userData)
	}
	ForwardRequest(w, r, getUserMgmtBaseURL(h.Config)+"/users", modifyBody)
}

func (h *Handler) SigninHandler(w http.ResponseWriter, r *http.Request) {
	ForwardRequest(w, r, getUserMgmtBaseURL(h.Config)+"/signin", nil)
}

func (h *Handler) SignoutHandler(w http.ResponseWriter, r *http.Request) {
	ForwardRequest(w, r, getUserMgmtBaseURL(h.Config)+"/signout", nil)
}

// TODO: Remove this from gpsd-api-gateway, only temporary
var secretKey = []byte("secret key")

func VerifyToken(tokenString string) (bool, error) {
	var errStr error = fmt.Errorf("invalid token")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return false, errStr
	}

	if !token.Valid {
		return false, errStr
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, errStr
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return false, errStr
	}

	if time.Now().Unix() > int64(exp) {
		return false, errStr
	}

	return true, nil
}

func (h *Handler) VerifyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.WriteHeader(http.StatusUnauthorized)
		if e := json.NewEncoder(w).Encode(ErrorResponse{Error: "no token provided"}); e != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	token := ""
	if len(authHeader) > 7 && strings.HasPrefix(authHeader, "Bearer ") {
		token = authHeader[7:]
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		if e := json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid authorization format"}); e != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	valid, err := VerifyToken(token)
	if !valid || err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		if e := json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()}); e != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if e := json.NewEncoder(w).Encode(map[string]string{"message": "valid token"}); e != nil {

		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	ForwardRequest(w, r, getUserMgmtBaseURL(h.Config)+"/users", nil)
}

func (h *Handler) GetUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	ForwardRequest(w, r, getUserMgmtBaseURL(h.Config)+"/users"+id, nil)
}

func (h *Handler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	ForwardRequest(w, r, getUserMgmtBaseURL(h.Config)+"/users"+id, nil)
}

func (h *Handler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	ForwardRequest(w, r, getUserMgmtBaseURL(h.Config)+"/users"+id, nil)
}
