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
		Balancer:               &kafkago.LeastBytes{},
		RequiredAcks:           kafkago.RequireOne,
		BatchTimeout:           10 * time.Millisecond,
		BatchSize:              50,
		AllowAutoTopicCreation: true,
	}
	return &Producer{writer: w}
}

// Publish gửi 1 message lên Kafka vào topic cụ thể.
func (p *Producer) Publish(ctx context.Context, topic string, key string, event any) error {
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

	msg := p.buildMessage(topic, key, data, event)
	return p.writer.WriteMessages(ctx, msg)
}

// PublishBatch gửi nhiều message cùng lúc (atomic write) vào topic cụ thể.
func (p *Producer) PublishBatch(ctx context.Context, topic string, keys []string, events []any) error {
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
		msgs[i] = p.buildMessage(topic, keys[i], data, event)
	}
	return p.writer.WriteMessages(ctx, msgs...)
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

// buildMessage gắn headers tự động dựa trên interface event implement và gán Topic.
func (p *Producer) buildMessage(topic string, key string, data []byte, event any) kafkago.Message {
	msg := kafkago.Message{
		Topic: topic,
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

	log.Printf("📤 [Kafka] publishing to topic=%s key=%s headers=%v", topic, key, msg.Headers)
	return msg
}
