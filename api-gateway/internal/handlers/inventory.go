package handlers

import (
	"api-gateway/internal/api/proto/inventory"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type InventoryHandler struct {
	inventoryClient inventory.InventoryServiceClient
}

func NewInventoryHandler(inventoryClient inventory.InventoryServiceClient) *InventoryHandler {
	return &InventoryHandler{inventoryClient: inventoryClient}
}

func (h *InventoryHandler) AddItem(rw http.ResponseWriter, r *http.Request) {
	var req inventory.AddItemRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(rw, "Invalid request", http.StatusBadRequest)
		return
	}

	res, err := h.inventoryClient.AddItem(context.Background(), &req)
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

func (h *InventoryHandler) UpdateItem(rw http.ResponseWriter, r *http.Request) {
	var req inventory.UpdateItemRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(rw, "Invalid request", http.StatusBadRequest)
		return
	}

	res, err := h.inventoryClient.UpdateItem(context.Background(), &req)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	err = json.NewEncoder(rw).Encode(res)
	if err != nil {
		http.Error(rw, "Failed to encode register response", http.StatusInternalServerError)
	}
}

func (h *InventoryHandler) GetItem(rw http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	itemID := chi.URLParam(r, "itemID")

	req := inventory.GetItemRequest{
		Id:     itemID,
		UserId: userID,
	}

	res, err := h.inventoryClient.GetItem(context.Background(), &req)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	err = json.NewEncoder(rw).Encode(res)
	if err != nil {
		http.Error(rw, "Failed to encode register response", http.StatusInternalServerError)
	}
}

func (h *InventoryHandler) GetInventory(rw http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	req := inventory.GetInventoryRequest{UserId: userID}

	res, err := h.inventoryClient.GetInventory(context.Background(), &req)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	err = json.NewEncoder(rw).Encode(res)
	if err != nil {
		http.Error(rw, "Failed to encode register response", http.StatusInternalServerError)
	}
}

func (h *InventoryHandler) DeleteItem(rw http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	itemID := chi.URLParam(r, "itemID")

	req := inventory.DeleteItemRequest{
		Id:     itemID,
		UserId: userID,
	}

	res, err := h.inventoryClient.DeleteItem(context.Background(), &req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	err = json.NewEncoder(rw).Encode(res)
	if err != nil {
		http.Error(rw, "Failed to encode register response", http.StatusInternalServerError)
	}
}
