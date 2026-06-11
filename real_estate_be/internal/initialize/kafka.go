package initialize

import (
	"context"
	"log"

	"real_estate_be/internal/global"
	kafkaconsumer "real_estate_be/internal/kafka"
	model "real_estate_be/internal/models"
	"real_estate_be/internal/sse"
	"real_estate_be/pkg/kafka"

	"gorm.io/gorm"
)

// InitKafka khởi tạo Producer, SSE Hub và Kafka consumers.
func InitKafka() *kafka.Producer {
	// Khởi tạo SSE hub global
	global.SSEHub = sse.NewHub()

	// Nếu không có Kafka broker, bỏ qua consumer
	if len(global.Config.Kafka.Brokers) == 0 {
		log.Println("⚠️ [Kafka] no brokers configured, skipping consumers")
		return nil
	}

	producer := kafka.NewProducer()
	return producer
}

// StartKafkaConsumers khởi động EnrichConsumer và NotifyConsumer trong goroutine.
func StartKafkaConsumers(ctx context.Context, db *gorm.DB) {
	if len(global.Config.Kafka.Brokers) == 0 {
		return
	}

	// Dùng chung 1 producer
	producer := kafka.NewProducer()

	// EnrichConsumer
	go func() {
		enrich := kafkaconsumer.NewEnrichConsumer(db, producer)
		defer enrich.Close()
		enrich.Start(ctx)
	}()

	// NotifyConsumer
	go func() {
		notify := kafkaconsumer.NewNotifyConsumer(db)
		defer notify.Close()
		notify.Start(ctx)
	}()

	log.Println("✅ Kafka consumers started")
}

// MigrateDb tự động migrate các bảng.
func MigrateDb(db *gorm.DB) {
	if err := db.AutoMigrate(
		&model.User{},
		&model.RealEstateModel{},
		&model.Notification{},
	); err != nil {
		log.Fatalf("❌ DB migration failed: %v", err)
	}
	log.Println("✅ DB migration completed")
}
