package initialize

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"real_estate_be/internal/global"
	kafkaconsumer "real_estate_be/internal/kafka"
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

}

// ensureTopicExists chủ động tạo topic nếu chưa tồn tại.
//
// LƯU Ý: Kafka hiện đại (2.x+) tự forward CreateTopics từ broker bất kỳ tới controller,
// nên KHÔNG cần tự dial controller nữa. Cách cũ (conn.Controller() rồi Dial host/port đó)
// rất dễ chết vì endpoint controller trả về theo ADVERTISED LISTENER của client:
// nếu compose advertise PLAINTEXT://127.0.0.1:9092 mà BE chạy trong container khác thì
// "127.0.0.1" là chính BE ⇒ dial tcp 127.0.0.1:9092: connection refused.
//
// Ở đây chỉ dial broker đầu tiên (đúng listener của môi trường đang chạy) và thử lại vài lần
// vì lúc BE khởi động, Kafka có thể chưa sẵn sàng.
func ensureTopicExists(brokers []string, topic string, numPartitions int, replicationFactor int) {
	if len(brokers) == 0 {
		return
	}

	topicConfig := kafkago.TopicConfig{
		Topic:             topic,
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	}

	const maxAttempts = 5
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		conn, err := kafkago.Dial("tcp", brokers[0])
		if err != nil {
			lastErr = err
			log.Printf("⚠️ [Kafka-Admin] chưa kết nối được broker %s (lần %d/%d): %v",
				brokers[0], attempt, maxAttempts, err)
			time.Sleep(time.Duration(attempt) * 2 * time.Second)
			continue
		}

		err = conn.CreateTopics(topicConfig)
		conn.Close()
		if err == nil {
			log.Printf("✅ [Kafka-Admin] đã tạo topic '%s' (%d partition)", topic, numPartitions)
			return
		}

		// Topic đã tồn tại là trường hợp bình thường → coi như xong
		if errors.Is(err, kafkago.TopicAlreadyExists) ||
			strings.Contains(strings.ToLower(err.Error()), "already exists") {
			log.Printf("ℹ️ [Kafka-Admin] topic '%s' đã tồn tại", topic)
			return
		}

		lastErr = err
		log.Printf("⚠️ [Kafka-Admin] tạo topic '%s' thất bại (lần %d/%d): %v", topic, attempt, maxAttempts, err)
		time.Sleep(time.Duration(attempt) * 2 * time.Second)
	}

	log.Printf("⚠️ [Kafka-Admin] bỏ qua tạo topic '%s' sau %d lần thử: %v", topic, maxAttempts, lastErr)
}
