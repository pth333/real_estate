package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"real_estate_be/internal/global"

	kafkago "github.com/segmentio/kafka-go"
)

// Producer gửi message đến Kafka
type Producer struct {
	writer *kafkago.Writer
	mu     sync.Mutex
	closed bool
}

func NewProducer() (*Producer, error) {
	brokers := global.Config.Kafka.Brokers
	topic := global.Config.Kafka.Topics.RealEstateCrawled

	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers is empty")
	}
	if topic == "" {
		return nil, errors.New("kafka topic real_estate_crawled is empty")
	}

	return &Producer{
		writer: &kafkago.Writer{
			Addr:                   kafkago.TCP(brokers...),
			Topic:                  topic,
			AllowAutoTopicCreation: true,
			Balancer:               &kafkago.LeastBytes{},
			RequiredAcks:           kafkago.RequireOne,
			BatchTimeout:           10 * time.Millisecond,
			BatchSize:              50,
		},
	}, nil
}

// Publish gửi 1 event bất kỳ lên Kafka.
// key = Kafka message key (thường là SourceURL), event = struct bất kỳ có thể Marshal.
// Event phải có embedded BaseEvent để header tự động gắn.
func (p *Producer) Publish(ctx context.Context, key string, event any) error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return errors.New("kafka producer is closed")
	}
	p.mu.Unlock()

	if p.writer == nil {
		return errors.New("kafka producer is not initialized")
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	msg := p.buildMessage(key, data, event)
	return p.writer.WriteMessages(ctx, msg)
}

// PublishBatch gửi nhiều event cùng lúc.
func (p *Producer) PublishBatch(ctx context.Context, keys []string, events []any) error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return errors.New("kafka producer is closed")
	}
	p.mu.Unlock()

	if p.writer == nil {
		return errors.New("kafka producer is not initialized")
	}

	if len(keys) != len(events) {
		return errors.New("keys and events length mismatch")
	}

	messages := make([]kafkago.Message, 0, len(events))
	for i, event := range events {
		data, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("marshal event at %d: %w", i, err)
		}
		messages = append(messages, p.buildMessage(keys[i], data, event))
	}

	return p.writer.WriteMessages(ctx, messages...)
}

// buildMessage tạo kafkago.Message từ key, data và event headers
func (p *Producer) buildMessage(key string, data []byte, event any) kafkago.Message {
	msg := kafkago.Message{
		Key:   []byte(key),
		Value: data,
		Time:  time.Now(),
		Headers: []kafkago.Header{
			{Key: HeaderContentType, Value: []byte(ContentTypeJSON)},
		},
	}

	// Trích xuất headers từ BaseEvent nếu event có embedded
	if ev, ok := event.(interface{ GetEventType() string }); ok {
		msg.Headers = append(msg.Headers, kafkago.Header{Key: HeaderEventType, Value: []byte(ev.GetEventType())})
	}
	if ev, ok := event.(interface{ GetSource() string }); ok {
		msg.Headers = append(msg.Headers, kafkago.Header{Key: HeaderSource, Value: []byte(ev.GetSource())})
	}
	if ev, ok := event.(interface{ GetVersion() string }); ok {
		msg.Headers = append(msg.Headers, kafkago.Header{Key: HeaderVersion, Value: []byte(ev.GetVersion())})
	}

	return msg
}

// SetTopic đổi topic của producer (dùng khi publish nhiều topic khác nhau)
func (p *Producer) SetTopic(topic string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.writer.Topic = topic
}

// Close đóng producer an toàn (thread-safe)
func (p *Producer) Close() error {
	if p == nil {
		return nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return nil
	}
	p.closed = true

	if p.writer != nil {
		return p.writer.Close()
	}
	return nil
}
