package kafka

import (
	"context"
	"log"

	"real_estate_be/internal/global"
	model "real_estate_be/internal/models"
	"real_estate_be/internal/repo"
	kafkapkg "real_estate_be/pkg/kafka"

	kafkago "github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

// EnrichConsumer xử lý RealEstateCrawledEvent, enrich và publish lên topic enriched.
type EnrichConsumer struct {
	consumer *kafkapkg.Consumer
	producer *kafkapkg.Producer
	repo     repo.RealEstateRepository
}

func NewEnrichConsumer(db *gorm.DB, producer *kafkapkg.Producer) *EnrichConsumer {
	cfg := global.Config.Kafka
	consumer := kafkapkg.NewConsumer(kafkapkg.ConsumerConfig{
		Topic:       cfg.Topics.RealEstateCrawled,
		GroupSuffix: "enrich",
		Handler:     nil, // set bên dưới
		Concurrency: 2,
	})

	ec := &EnrichConsumer{
		consumer: consumer,
		producer: producer,
		repo:     repo.NewRealEstateRepository(db),
	}

	// Tạo consumer mới với handler đã biết
	ec.consumer = kafkapkg.NewConsumer(kafkapkg.ConsumerConfig{
		Topic:       cfg.Topics.RealEstateCrawled,
		GroupSuffix: "enrich",
		Handler:     ec.handle,
		Concurrency: 2,
	})

	return ec
}

func (e *EnrichConsumer) Start(ctx context.Context) {
	log.Println("🚀 [EnrichConsumer] starting...")

	// Set producer topic sang enriched topic
	e.producer.SetTopic(global.Config.Kafka.Topics.RealEstateEnriched)

	e.consumer.Start(ctx)
}

func (e *EnrichConsumer) Close() error {
	return e.consumer.Close()
}

func (e *EnrichConsumer) handle(ctx context.Context, msg kafkago.Message) error {
	// Parse event
	event, err := kafkapkg.UnmarshalMsg[kafkapkg.RealEstateCrawledEvent](msg)
	if err != nil {
		log.Printf("❌ [EnrichConsumer] unmarshal error: %v", err)
		return nil // không retry
	}

	log.Printf("🔧 [EnrichConsumer] processing: %s", event.SourceURL)

	// Phân loại BĐS — hiện là stub
	typeStr := classifyType(event.Title, event.PriceVND)

	// Lưu DB với enriched data
	enriched := &model.RealEstate{
		Title:      event.Title,
		PriceVND:   event.PriceVND,
		Address:    event.Address,
		District:   event.District,
		City:       event.City,
		Acreage:    event.Acreage,
		PricePerM2: event.PricePerM2,
	}
	if err := e.repo.Create(enriched); err != nil {
		log.Printf("⚠️ [EnrichConsumer] DB error: %v", err)
		// Lỗi DB vẫn tiếp tục — không block
	}

	// Publish enriched event
	enrichedEvent := kafkapkg.NewRealEstateEnrichedEvent(*event, typeStr)
	if err := e.producer.Publish(ctx, event.SourceURL, enrichedEvent); err != nil {
		log.Printf("⚠️ [EnrichConsumer] publish error: %v", err)
	}

	return nil
}

// classifyType — stub, luôn trả "apartment".
func classifyType(title string, priceVND float64) string {
	// TODO: dùng title + price để phân loại chính xác hơn
	return "apartment"
}
