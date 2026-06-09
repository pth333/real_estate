package kafka

import (
	"context"
	"log"

	"real_estate_be/internal/global"
	"real_estate_be/pkg/kafka"

	kafkago "github.com/segmentio/kafka-go"
)

// NotifyConsumer nhận enriched events và gửi notification.
// Hiện tại là stub — sẽ implement sau.
type NotifyConsumer struct{}

func NewNotifyConsumer() *NotifyConsumer {
	return &NotifyConsumer{}
}

func (n *NotifyConsumer) Start(ctx context.Context) error {
	consumer := kafka.NewConsumer(kafka.ConsumerConfig{
		Topic:       global.Config.Kafka.Topics.RealEstateEnriched,
		GroupSuffix: "notify",
		Handler:     n.handle,
		Concurrency: 1,
	})

	return consumer.Start(ctx)
}

func (n *NotifyConsumer) handle(ctx context.Context, msg kafkago.Message) error {
	eventType := kafka.GetEventHeader(msg, kafka.HeaderEventType)
	if eventType != kafka.EventTypeEnriched {
		log.Printf("[Notify] Skip unknown event type: %s", eventType)
		return nil
	}

	var event kafka.RealEstateEnrichedEvent
	if err := kafka.UnmarshalEvent(msg, &event); err != nil {
		log.Printf("[Notify] Parse error: %v", err)
		return nil
	}

	log.Printf("[Notify] Received enriched: %s — %s (type=%s)",
		event.SourceURL, event.Title, event.TypeOfRealEstate)

	// TODO: gửi email/SMS/webhook cho user đã subscribe
	// TODO: publish RealEstateNotifiedEvent sau khi gửi thành công

	return nil
}
