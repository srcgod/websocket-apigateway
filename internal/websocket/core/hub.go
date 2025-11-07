package core

import (
	"sync"

	"github.com/sirupsen/logrus"
)

// Hub управляет всеми активными WebSocket соединениями
type Hub struct {
	clients    map[string]*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	logger     *logrus.Logger

	mu sync.RWMutex
}

// NewHub создает новый Hub
func NewHub(logger *logrus.Logger) *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		logger:     logger,
	}
}

// HubRun запускает основной цикл Hub
func (h *Hub) HubRun() {
	h.logger.Info("WebSocket Hub started!")
	for {
		select {
		case client := <-h.register:
			h.handleRegister(client)

		case client := <-h.unregister:
			h.handleUnregister(client)
		}
	}
}

// RegisterClient регистрирует нового клиента
func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

// UnregisterClient удаляет клиента
func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

// handleRegister обрабатывает регистрацию клиента
func (h *Hub) handleRegister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.clients[client.ID()]; exists {
		h.logger.Warnf("Client %s already registered", client.ID())
		return
	}

	h.clients[client.ID()] = client

	h.logger.WithFields(logrus.Fields{
		"client_id": client.ID(),
		"user_id":   client.UserID(),
		"total":     len(h.clients),
	}).Info("🔌 Client registered")
}

// handleUnregister обрабатывает удаление клиента
func (h *Hub) handleUnregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if storedClient, exists := h.clients[client.ID()]; exists {
		close(storedClient.Send)
		delete(h.clients, client.ID())

		h.logger.WithFields(logrus.Fields{
			"client_id": client.ID(),
			"user_id":   client.UserID(),
			"total":     len(h.clients),
		}).Info("🔌 Client unregistered")
	}
}

// TODO: реализовать broadcast
func (h *Hub) broadcastTo() {
	// impl...
}

