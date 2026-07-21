package routers

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"time"

	"real_estate_be/internal/global"
	"real_estate_be/internal/sse"

	"github.com/gofiber/fiber/v2"
)

func InitNotificationRoutes(Router fiber.Router) {
	notifRouter := Router.Group("/notifications")
	{
		notifRouter.Get("/stream", SSEStream)
		notifRouter.Get("/", ListNotifications)
		notifRouter.Patch("/:id/read", MarkAsRead)
	}
}

// SSEStream trả về Server-Sent Events stream cho user hiện tại.
func SSEStream(c *fiber.Ctx) error {
	userIDStr := c.Query("user_id", "0")
	userID, _ := strconv.ParseUint(userIDStr, 10, 64)

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Access-Control-Allow-Origin", "*")

	client := &sse.Client{
		UserID: userID,
		Send:   make(chan []byte, 64),
		Done:   make(chan struct{}),
	}

	if global.SSEHub != nil {
		global.SSEHub.Register(client)
		defer global.SSEHub.Unregister(client)
	}

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		defer close(client.Done)

		// keepalive ticker
		keepaliveTicker := time.NewTicker(30 * time.Second)
		defer keepaliveTicker.Stop()

		// Luồng đọc context.Done
		done := make(chan struct{})
		go func() {
			ctx := c.Context()
			if ctx != nil {
				<-ctx.Done()
			}
			close(done)
		}()

		for {
			select {
			case <-done:
				return
			case data, ok := <-client.Send:
				if !ok {
					return
				}
				msg := fmt.Sprintf("data: %s\n\n", string(data))
				if _, err := w.WriteString(msg); err != nil {
					log.Printf("⚠️ [SSE] write error: %v", err)
					return
				}
				if err := w.Flush(); err != nil {
					log.Printf("⚠️ [SSE] flush error: %v", err)
					return
				}
			case <-keepaliveTicker.C:
				if _, err := w.WriteString(": keepalive\n\n"); err != nil {
					return
				}
				if err := w.Flush(); err != nil {
					return
				}
			}
		}
	})

	return nil
}

// ListNotifications trả về danh sách notification của user (phân trang).
func ListNotifications(c *fiber.Ctx) error {
	userIDStr := c.Query("user_id", "0")
	userID, _ := strconv.ParseUint(userIDStr, 10, 64)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	var notifications []map[string]interface{}
	if err := global.DB.Table("notifications").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&notifications).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	var total int64
	global.DB.Model(&struct{}{}).Table("notifications").
		Where("user_id = ?", userID).
		Count(&total)

	return c.JSON(fiber.Map{
		"data":  notifications,
		"total": total,
		"page":  page,
	})
}

// MarkAsRead đánh dấu 1 notification đã đọc.
func MarkAsRead(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	if err := global.DB.Table("notifications").Where("id = ?", id).
		Update("is_read", true).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}
