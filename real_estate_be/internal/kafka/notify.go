package kafka

import (
	"context"
	"encoding/json"
	"log"

	"real_estate_be/internal/global"
	model "real_estate_be/internal/models"
	kafkapkg "real_estate_be/pkg/kafka"

	kafkago "github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

// NotifyConsumer xử lý RealEstateEnrichedEvent, tạo notification và broadcast SSE.
type NotifyConsumer struct {
	consumer *kafkapkg.Consumer
	db       *gorm.DB
}

func NewNotifyConsumer(db *gorm.DB) *NotifyConsumer {
	cfg := global.Config.Kafka
	consumer := kafkapkg.NewConsumer(kafkapkg.ConsumerConfig{
		Topic:       cfg.Topics.RealEstateEnriched,
		GroupSuffix: "notify",
		Handler:     nil,
		Concurrency: 1,
	})

	nc := &NotifyConsumer{
		consumer: consumer,
		db:       db,
	}

	nc.consumer = kafkapkg.NewConsumer(kafkapkg.ConsumerConfig{
		Topic:       cfg.Topics.RealEstateEnriched,
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
	event, err := kafkapkg.UnmarshalMsg[kafkapkg.RealEstateEnrichedEvent](msg)
	if err != nil {
		log.Printf("❌ [NotifyConsumer] unmarshal error: %v", err)
		return nil
	}

	log.Printf("🔔 [NotifyConsumer] new listing: %s", event.SourceURL)

	// Lưu notification vào DB cho tất cả user
	// Broadcast message
	type notificationPayload struct {
		Type    string `json:"type"`
		Title   string `json:"title"`
		Message string `json:"message"`
		Price   float64 `json:"price_vnd"`
		Area    float64 `json:"acreage"`
		Address string  `json:"address"`
		URL     string  `json:"source_url"`
	}

	payload := notificationPayload{
		Type:    "new_listing",
		Title:   event.Title,
		Message: "Bất động sản mới: " + event.Title,
		Price:   event.PriceVND,
		Area:    event.Acreage,
		Address: event.Address,
		URL:     event.SourceURL,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Lưu notification cho mỗi user vào DB
	n.saveNotifications(event)

	// Broadcast SSE đến tất cả client đang online
	hub := global.SSEHub
	if hub != nil {
		hub.Broadcast(data)
	}

	return nil
}

// saveNotifications lưu notification record cho mỗi user.
func (n *NotifyConsumer) saveNotifications(event *kafkapkg.RealEstateEnrichedEvent) {
	var users []model.User
	if err := n.db.Model(&model.User{}).Find(&users).Error; err != nil {
		log.Printf("⚠️ [NotifyConsumer] get users error: %v", err)
		return
	}

	for _, user := range users {
		notif := model.Notification{
			UserID:  user.ID,
			Title:   "Bất động sản mới",
			Message: event.Title + " - " + event.Address,
		}
		if err := n.db.Create(&notif).Error; err != nil {
			log.Printf("⚠️ [NotifyConsumer] save notification error: %v", err)
		}
	}
}
