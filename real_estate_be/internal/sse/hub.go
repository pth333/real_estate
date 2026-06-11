package sse

import (
	"log"
	"sync"
)

// Client đại diện cho một kết nối SSE.
type Client struct {
	UserID uint64
	Send   chan []byte
	Done   chan struct{}
}

// Hub quản lý tất cả client đang kết nối.
type Hub struct {
	mu      sync.RWMutex
	clients map[uint64]map[*Client]bool
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[uint64]map[*Client]bool),
	}
}

// Register thêm client vào hub.
func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[client.UserID] == nil {
		h.clients[client.UserID] = make(map[*Client]bool)
	}
	h.clients[client.UserID][client] = true
	log.Printf("🔌 [SSE] client connected: userID=%d (total connections: %d)", client.UserID, h.count())
}

// Unregister xoá client khỏi hub và đóng channel.
func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.clients[client.UserID]; ok {
		if _, exists := clients[client]; exists {
			delete(clients, client)
			close(client.Send)
			if len(clients) == 0 {
				delete(h.clients, client.UserID)
			}
		}
	}
	log.Printf("🔌 [SSE] client disconnected: userID=%d", client.UserID)
}

// Broadcast gửi message đến tất cả client.
func (h *Hub) Broadcast(data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, clients := range h.clients {
		for client := range clients {
			select {
			case client.Send <- data:
			default:
				// Bỏ qua nếu buffer đầy (client chậm)
			}
		}
	}
}

// BroadcastToUser gửi message đến 1 user cụ thể.
func (h *Hub) BroadcastToUser(userID uint64, data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.clients[userID]; ok {
		for client := range clients {
			select {
			case client.Send <- data:
			default:
			}
		}
	}
}

func (h *Hub) count() int {
	total := 0
	for _, clients := range h.clients {
		total += len(clients)
	}
	return total
}
