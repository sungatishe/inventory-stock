package handlers

import (
	"api-gateway/internal/api/proto/auth"
	"context"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	authClient auth.AuthServiceClient
}

func NewAuthHandler(authClient auth.AuthServiceClient) *AuthHandler {
	return &AuthHandler{authClient: authClient}
}

func (h *AuthHandler) Register(rw http.ResponseWriter, r *http.Request) {
	var req auth.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(rw, "Invalid request", http.StatusBadRequest)
		return
	}

	res, err := h.authClient.Register(context.Background(), &req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(rw).Encode(res)
	if err != nil {
		http.Error(rw, "Failed to encode register response", http.StatusInternalServerError)
	}
}

func (h *AuthHandler) Login(rw http.ResponseWriter, r *http.Request) {
	var req auth.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := h.authClient.Login(context.Background(), &req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	err = json.NewEncoder(rw).Encode(res)
	if err != nil {
		http.Error(rw, "Failed to encode login response", http.StatusInternalServerError)
	}
}
