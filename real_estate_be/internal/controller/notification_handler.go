package controller

import (
	"bufio"
	"encoding/json"
	"fmt"
	"real_estate_be/internal/global"
	"real_estate_be/internal/response"
	"real_estate_be/internal/sse"
	"real_estate_be/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	service usecase.INotificationService
}

func NewNotificationHandler(service usecase.INotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

// Stream xử lý kết nối SSE
func (h *NotificationHandler) Stream(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")
	c.Set("X-Accel-Buffering", "no") // Quan trọng khi chạy sau Nginx

	id := uuid.New().String()
	client := &sse.Client{
		ID:   id,
		Send: make(chan []byte, 10),
	}

	global.SSEHub.Register(client)

	// Lấy context từ request để theo dõi việc đóng kết nối an toàn hơn
	ctx := c.UserContext()

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		fmt.Fprintf(w, "retry: 5000\n\n")
		w.Flush()

		for {
			select {
			case msg, ok := <-client.Send:
				if !ok {
					return
				}
				fmt.Fprintf(w, "data: %s\n\n", msg)
				if err := w.Flush(); err != nil {
					global.SSEHub.Unregister(id)
					return
				}
			case <-ctx.Done():
				global.SSEHub.Unregister(id)
				return
			}
		}
	})

	return nil
}

// GetNotifications lấy danh sách thông báo cũ
func (h *NotificationHandler) GetNotifications(c *fiber.Ctx) error {
	notifications, err := h.service.GetNotifications()
	if err != nil {
		return response.Error(c, 500, err.Error(), nil)
	}

	type notifResponse struct {
		ID        uint64      `json:"id"`
		ListingID uint64      `json:"listing_id"`
		Type      string      `json:"type"`
		Payload   interface{} `json:"payload"`
		CreatedAt string      `json:"created_at"`
	}

	var results []notifResponse
	for _, n := range notifications {
		var payload map[string]interface{}
		json.Unmarshal([]byte(n.Payload), &payload)
		results = append(results, notifResponse{
			ID:        n.ID,
			ListingID: n.ListingID,
			Type:      n.Type,
			Payload:   payload,
			CreatedAt: n.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return response.OK(c, results)
}
