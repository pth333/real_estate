package kafka

import (
	"context"
	"log"
	"time"

	"real_estate_be/internal/global"
	model "real_estate_be/internal/models"
	"real_estate_be/pkg/kafka"

	kafkago "github.com/segmentio/kafka-go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// EnrichConsumer nhận crawled events, enrich và lưu vào DB + publish enriched event
type EnrichConsumer struct {
	db       *gorm.DB
	producer *kafka.Producer
}

func NewEnrichConsumer(db *gorm.DB) *EnrichConsumer {
	return &EnrichConsumer{db: db}
}

func (e *EnrichConsumer) Start(ctx context.Context) error {
	producer, err := kafka.NewProducer()
	if err != nil {
		return err
	}
	e.producer = producer
	defer e.producer.Close()

	// Override topic của producer sang enriched topic
	e.producer.SetTopic(global.Config.Kafka.Topics.RealEstateEnriched)

	consumer := kafka.NewConsumer(kafka.ConsumerConfig{
		Topic:       global.Config.Kafka.Topics.RealEstateCrawled,
		GroupSuffix: "enrich",
		Handler:     e.handle,
		Concurrency: 2,
	})

	return consumer.Start(ctx)
}

func (e *EnrichConsumer) handle(ctx context.Context, msg kafkago.Message) error {
	eventType := kafka.GetEventHeader(msg, kafka.HeaderEventType)
	if eventType != kafka.EventTypeCrawled {
		log.Printf("[Enrich] Skip unknown event type: %s", eventType)
		return nil
	}

	var event kafka.RealEstateCrawledEvent
	if err := kafka.UnmarshalEvent(msg, &event); err != nil {
		log.Printf("[Enrich] Parse error: %v", err)
		return nil
	}

	log.Printf("[Enrich] Processing: %s — %s", event.SourceURL, event.Title)

	// ---- Enrichment ----
	typeOfRealEstate := classifyType(event.Title, event.PriceVND)
	var lat, lng *float64

	// Lưu enrichment vào DB
	enriched := model.RealEstateEnriched{
		SourceURL:        event.SourceURL,
		TypeOfRealEstate: typeOfRealEstate,
		Latitude:         lat,
		Longitude:        lng,
	}

	if err := e.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "source_url"}},
		DoUpdates: clause.AssignmentColumns([]string{"type_of_real_estate", "latitude", "longitude"}),
	}).Create(&enriched).Error; err != nil {
		log.Printf("[Enrich] DB error: %v", err)
		return nil
	}

	// Publish enriched event
	enrichedEvent := kafka.RealEstateEnrichedEvent{
		BaseEvent: kafka.BaseEvent{
			EventType: kafka.EventTypeEnriched,
			Source:    event.Source,
			Version:   "1.0",
			Timestamp: time.Now(),
		},
		SourceURL:        event.SourceURL,
		Title:            event.Title,
		Address:          event.Address,
		District:         event.District,
		City:             event.City,
		PriceVND:         event.PriceVND,
		Acreage:          event.Acreage,
		PricePerM2:       event.PricePerM2,
		TypeOfRealEstate: typeOfRealEstate,
		Latitude:         lat,
		Longitude:        lng,
	}

	if err := e.producer.Publish(ctx, event.SourceURL, enrichedEvent); err != nil {
		log.Printf("[Enrich] Publish error: %v", err)
		return nil
	}

	log.Printf("[Enrich] Done: %s → %s", event.SourceURL, typeOfRealEstate)
	return nil
}

// classifyType heuristic phân loại BĐS dựa trên title
func classifyType(title string, price float64) string {
	// TODO: có thể train model sau này, giờ rule-based tạm
	return "apartment"
}
