package initialize

import (
	"context"
	"log"

	"real_estate_be/internal/global"
	kafkaconsumer "real_estate_be/internal/kafka"
)

// StartKafkaConsumers khởi động các Kafka consumer workers.
// Hàm block nếu consumer Start() block;
// gọi trong goroutine nếu cần non-blocking.
func StartKafkaConsumers(ctx context.Context) {
	if len(global.Config.Kafka.Brokers) == 0 {
		log.Println("[Kafka] No brokers configured — skip consumers")
		return
	}

	// Enrich consumer: lấy crawled events → enrich → publish enriched
	enrich := kafkaconsumer.NewEnrichConsumer(global.DB)
	go func() {
		log.Println("[Kafka] Starting EnrichConsumer...")
		if err := enrich.Start(ctx); err != nil {
			log.Printf("[Kafka] EnrichConsumer stopped: %v", err)
		}
	}()

	// Notify consumer (stub): lấy enriched events → gửi notification (future)
	notify := kafkaconsumer.NewNotifyConsumer()
	go func() {
		log.Println("[Kafka] Starting NotifyConsumer...")
		if err := notify.Start(ctx); err != nil {
			log.Printf("[Kafka] NotifyConsumer stopped: %v", err)
		}
	}()

	log.Println("[Kafka] Consumers started ✅")
}
