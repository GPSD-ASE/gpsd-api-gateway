package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

var client = &http.Client{
	Timeout: 10 * time.Second,
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func ForwardRequest(w http.ResponseWriter, r *http.Request, endpoint string, modifyBody func([]byte) ([]byte, error)) {
	w.Header().Set("Content-Type", "application/json")

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

	newReq, err := http.NewRequest(r.Method, endpoint, bytes.NewBuffer(actualBody))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	newReq.Header.Set("Content-Type", "application/json")
	if auth := r.Header.Get("Authorization"); auth != "" {
		newReq.Header.Set("Authorization", auth)
	}

	for name, values := range r.Header {
		if name != "Content-Length" && name != "Transfer-Encoding" && name != "Connection" {
			for _, value := range values {
				newReq.Header.Add(name, value)
			}
		}
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
