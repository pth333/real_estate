package sse

import (
	"log"
	"sync"
)

// Client đại diện cho một kết nối SSE.
type Client struct {
	ID   string
	Send chan []byte
}

// Hub quản lý tất cả client đang kết nối.
type Hub struct {
	mu      sync.RWMutex
	clients map[string]*Client
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*Client),
	}
}

// Register thêm client vào hub.
func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client.ID] = client
	log.Printf("🔌 [SSE] client connected: id=%s (total connections: %d)", client.ID, len(h.clients))
}

// Unregister xoá client khỏi hub và đóng channel.
func (h *Hub) Unregister(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if client, ok := h.clients[id]; ok {
		delete(h.clients, id)
		close(client.Send)
		log.Printf("🔌 [SSE] client disconnected: id=%s", id)
	}
}

// Broadcast gửi message đến tất cả client.
func (h *Hub) Broadcast(data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.clients {
		select {
		case client.Send <- data:
		default:
			// Bỏ qua nếu buffer đầy (client chậm)
		}
	}
}
