package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"real_estate_be/internal/global"
	model "real_estate_be/internal/models"
	"real_estate_be/internal/repo"
	kafkapkg "real_estate_be/pkg/kafka"

	kafkago "github.com/segmentio/kafka-go"
)

// NotifyConsumer xử lý RealEstateNewListingEvent, tạo notification và broadcast SSE.
type NotifyConsumer struct {
	consumer *kafkapkg.Consumer
	repo     repo.INotificationRepository
}

func NewNotifyConsumer(repo repo.INotificationRepository) *NotifyConsumer {
	cfg := global.Config.Kafka

	topic := cfg.Topics.RealEstateNotified
	if topic == "" {
		topic = "real_estate.notified.v1"
	}

	nc := &NotifyConsumer{
		repo: repo,
	}

	nc.consumer = kafkapkg.NewConsumer(kafkapkg.ConsumerConfig{
		Topic:       topic,
		GroupSuffix: "notify",
		Handler:     nc.handle,
		Concurrency: 1,
	})

	return nc
}

func (n *NotifyConsumer) Start(ctx context.Context) {
	log.Println("🔔 [NotifyConsumer] starting...")
	n.consumer.Start(ctx)
}

func (n *NotifyConsumer) Close() error {
	return n.consumer.Close()
}

func (n *NotifyConsumer) handle(ctx context.Context, msg kafkago.Message) error {
	event, err := kafkapkg.UnmarshalMsg[kafkapkg.RealEstateNewListingEvent](msg)
	if err != nil {
		log.Printf("❌ [NotifyConsumer] unmarshal error: %v", err)
		return nil
	}

	log.Printf("🔔 [NotifyConsumer] new notify for listing: %d", event.ListingID)

	// Tạo payload cho SSE và DB
	payloadMap := map[string]interface{}{
		"title":   event.Title,
		"address": event.Address,
		"price":   event.PriceVND,
		"acreage": event.Acreage,
		"slug":    event.Slug,
	}

	payloadBytes, _ := json.Marshal(payloadMap)

	// Lưu notification vào DB
	notif := &model.Notification{
		ListingID: event.ListingID,
		Type:      "new_listing",
		Payload:   string(payloadBytes),
		CreatedAt: time.Now(),
	}

	if err := n.repo.Create(notif); err != nil {
		log.Printf("⚠️ [NotifyConsumer] save notification error: %v", err)
	}

	// Broadcast SSE
	if global.SSEHub != nil {
		global.SSEHub.Broadcast(payloadBytes)
	}

	return nil
}
