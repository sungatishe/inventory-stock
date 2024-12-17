package handlers

import (
	"api-gateway/internal/ws"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

type StockHandlers struct {
	manager *ws.Manager
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewStockHandlers(manager *ws.Manager) *StockHandlers {
	return &StockHandlers{manager: manager}
}

func (h *StockHandlers) NotificationHandler(rw http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	if userIDStr == "" {
		http.Error(rw, "userID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(rw, "invalid userID", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v\n", err)
		return
	}

	h.manager.AddClient(userID, conn)

	defer h.manager.RemoveClient(userID)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from client %d: %v\n", userID, err)
			break
		}
	}
}
