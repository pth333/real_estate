package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"real_estate_be/internal/global"

	kafkago "github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafkago.Writer
	mu     sync.Mutex
	closed bool
}

func NewProducer() *Producer {
	cfg := global.Config.Kafka
	w := &kafkago.Writer{
		Addr:                   kafkago.TCP(cfg.Brokers...),
		Topic:                  cfg.Topics.RealEstateNotified,
		Balancer:               &kafkago.LeastBytes{},
		RequiredAcks:           kafkago.RequireOne,
		BatchTimeout:           10 * time.Millisecond,
		BatchSize:              50,
		AllowAutoTopicCreation: true,
	}
	return &Producer{writer: w}
}

// Publish gửi 1 message lên Kafka.
func (p *Producer) Publish(ctx context.Context, key string, event any) error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return fmt.Errorf("producer is closed")
	}
	p.mu.Unlock()

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	msg := p.buildMessage(key, data, event)
	return p.writer.WriteMessages(ctx, msg)
}

// PublishBatch gửi nhiều message cùng lúc (atomic write).
func (p *Producer) PublishBatch(ctx context.Context, keys []string, events []any) error {
	if len(keys) != len(events) {
		return fmt.Errorf("keys and events length mismatch")
	}
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return fmt.Errorf("producer is closed")
	}
	p.mu.Unlock()

	msgs := make([]kafkago.Message, len(events))
	for i, event := range events {
		data, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("json.Marshal[%d]: %w", i, err)
		}
		msgs[i] = p.buildMessage(keys[i], data, event)
	}
	return p.writer.WriteMessages(ctx, msgs...)
}

// SetTopic cho phép đổi topic — dùng trong EnrichConsumer.
func (p *Producer) SetTopic(topic string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.writer.Topic = topic
}

func (p *Producer) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return nil
	}
	p.closed = true
	return p.writer.Close()
}

// buildMessage gắn headers tự động dựa trên interface event implement.
func (p *Producer) buildMessage(key string, data []byte, event any) kafkago.Message {
	msg := kafkago.Message{
		Key:   []byte(key),
		Value: data,
		Headers: []kafkago.Header{
			{Key: HeaderContentType, Value: []byte("application/json")},
		},
	}

	if e, ok := event.(EventTyper); ok {
		msg.Headers = append(msg.Headers, kafkago.Header{Key: HeaderEventType, Value: []byte(e.GetEventType())})
	}
	if e, ok := event.(EventSourcer); ok {
		msg.Headers = append(msg.Headers, kafkago.Header{Key: HeaderSource, Value: []byte(e.GetSource())})
	}
	if e, ok := event.(EventVersion); ok {
		msg.Headers = append(msg.Headers, kafkago.Header{Key: HeaderVersion, Value: []byte(e.GetVersion())})
	}

	log.Printf("📤 [Kafka] publishing key=%s headers=%v", key, msg.Headers)
	return msg
}
