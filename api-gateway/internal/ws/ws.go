package ws

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type Client struct {
	ID   int
	Conn *websocket.Conn
}

type Manager struct {
	Clients map[int]*Client
	mu      sync.Mutex
}

func NewManager() *Manager {
	return &Manager{Clients: make(map[int]*Client)}
}

func (m *Manager) AddClient(userID int, conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Clients[userID] = &Client{ID: userID, Conn: conn}
	log.Printf("Client %d connected\n", userID)
}

func (m *Manager) RemoveClient(userID int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if client, ok := m.Clients[userID]; ok {
		client.Conn.Close()
		delete(m.Clients, userID)
		log.Printf("Client %d disconnected\n", userID)
	}
}

func (m *Manager) SendNotification(userID int, message string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	client, ok := m.Clients[userID]
	if !ok {
		log.Printf("Client %d not connected. Skipping notification\n", userID)
		return
	}

	err := client.Conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Printf("Failed to send message to client %d: %v\n", userID, err)
		m.RemoveClient(userID)
	}
}
