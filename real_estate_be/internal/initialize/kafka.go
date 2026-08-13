package initialize

import (
	"context"
	"log"
	"net"
	"strconv"

	"real_estate_be/internal/global"
	kafkaconsumer "real_estate_be/internal/kafka"
	model "real_estate_be/internal/models"
	"real_estate_be/internal/repo"
	"real_estate_be/internal/sse"
	"real_estate_be/pkg/kafka"

	kafkago "github.com/segmentio/kafka-go"
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

	// 1. Đảm bảo các topic cần thiết đã tồn tại trước khi khởi chạy consumer
	topicsToCreate := []string{
		global.Config.Kafka.Topics.RealEstateNotified,
	}
	for _, t := range topicsToCreate {
		if t != "" {
			ensureTopicExists(global.Config.Kafka.Brokers, t, 1, 1)
		}
	}

	// Dùng chung 1 producer
	// producer := kafka.NewProducer()

	// EnrichConsumer (Tạm thời bị comment ở repo)
	// go func() {
	// 	enrich := kafkaconsumer.NewEnrichConsumer(db, producer)
	// 	defer enrich.Close()
	// 	enrich.Start(ctx)
	// }()

	// NotifyConsumer
	go func() {
		notificationRepo := repo.NewNotificationRepository(db)
		notify := kafkaconsumer.NewNotifyConsumer(notificationRepo)
		defer notify.Close()
		notify.Start(ctx)
	}()

	log.Println("✅ Kafka consumers started")
}

// ensureTopicExists chủ động tạo topic nếu chưa tồn tại
func ensureTopicExists(brokers []string, topic string, numPartitions int, replicationFactor int) {
	if len(brokers) == 0 {
		return
	}

	// Kết nối tới broker đầu tiên để gửi yêu cầu quản trị (Admin)
	conn, err := kafkago.Dial("tcp", brokers[0])
	if err != nil {
		log.Printf("⚠️ [Kafka-Admin] failed to connect to broker %s: %v", brokers[0], err)
		return
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		log.Printf("⚠️ [Kafka-Admin] failed to get controller: %v", err)
		return
	}

	var controllerConn *kafkago.Conn
	controllerConn, err = kafkago.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		log.Printf("⚠️ [Kafka-Admin] failed to connect to controller: %v", err)
		return
	}
	defer controllerConn.Close()

	topicConfig := kafkago.TopicConfig{
		Topic:             topic,
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	}

	err = controllerConn.CreateTopics(topicConfig)
	if err != nil {
		// Nếu topic đã tồn tại thì broker sẽ trả về lỗi, chúng ta có thể an tâm bỏ qua
		log.Printf("ℹ️ [Kafka-Admin] Topic '%s' check completed (it may already exist or auto-created: %v)", topic, err)
		return
	}

	log.Printf("🎉 [Kafka-Admin] Topic '%s' created successfully with %d partitions", topic, numPartitions)
}

// MigrateDb tự động migrate các bảng.
func MigrateDb(db *gorm.DB) {
	if err := db.AutoMigrate(
		&model.User{},
		&model.RealEstate{},
		&model.Notification{},
	); err != nil {
		log.Fatalf("❌ DB migration failed: %v", err)
	}
	log.Println("✅ DB migration completed")
}
