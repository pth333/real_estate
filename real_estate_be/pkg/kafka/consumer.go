package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"real_estate_be/internal/global"

	kafkago "github.com/segmentio/kafka-go"
)

// ConsumerHandler xử lý từng message nhận được
type ConsumerHandler func(ctx context.Context, msg kafkago.Message) error

// Consumer là base consumer reusable
type Consumer struct {
	reader      *kafkago.Reader
	topic       string
	groupID     string
	handler     ConsumerHandler
	concurrency int
}

type ConsumerConfig struct {
	Topic       string
	GroupSuffix string // suffix thêm vào group_prefix từ config
	Handler     ConsumerHandler
	Concurrency int // số goroutine xử lý song song (default 1)
}

func NewConsumer(cfg ConsumerConfig) *Consumer {
	brokers := global.Config.Kafka.Brokers
	groupPrefix := global.Config.Kafka.GroupPrefix

	groupID := groupPrefix
	if cfg.GroupSuffix != "" {
		groupID = groupPrefix + "-" + cfg.GroupSuffix
	}

	if cfg.Concurrency < 1 {
		cfg.Concurrency = 1
	}

	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:        brokers,
		Topic:          cfg.Topic,
		GroupID:        groupID,
		GroupTopics:    []string{cfg.Topic},
		StartOffset:    kafkago.FirstOffset,
		MinBytes:       10,
		MaxBytes:       10e6, // 10MB
		MaxWait:        1 * time.Second,
		CommitInterval: time.Second,
		RetentionTime:  24 * time.Hour,
	})

	return &Consumer{
		reader:      reader,
		topic:       cfg.Topic,
		groupID:     groupID,
		handler:     cfg.Handler,
		concurrency: cfg.Concurrency,
	}
}

// Start bắt đầu consume message, block đến khi context bị cancel
func (c *Consumer) Start(ctx context.Context) error {
	log.Printf("[Kafka] Consumer started — topic=%s group=%s concurrency=%d",
		c.topic, c.groupID, c.concurrency)

	if c.concurrency <= 1 {
		return c.loop(ctx)
	}

	return c.loopConcurrent(ctx)
}

// loop xử lý tuần tự
func (c *Consumer) loop(ctx context.Context) error {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			return fmt.Errorf("fetch message: %w", err)
		}

		if err := c.handler(ctx, msg); err != nil {
			log.Printf("[Kafka] Handler error: %v — msg key=%s", err, string(msg.Key))
			// Commit anyway để không block consumer
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			return fmt.Errorf("commit message: %w", err)
		}
	}
}

// loopConcurrent xử lý song song với semaphore
func (c *Consumer) loopConcurrent(ctx context.Context) error {
	sem := make(chan struct{}, c.concurrency)
	errCh := make(chan error, 1)

	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			return fmt.Errorf("fetch message: %w", err)
		}

		sem <- struct{}{}
		go func(m kafkago.Message) {
			defer func() { <-sem }()

			if err := c.handler(ctx, m); err != nil {
				log.Printf("[Kafka] Handler error: %v — msg key=%s", err, string(m.Key))
			}

			if err := c.reader.CommitMessages(ctx, m); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}(msg)

		select {
		case commitErr := <-errCh:
			return commitErr
		default:
		}
	}
}

// Close đóng consumer
func (c *Consumer) Close() error {
	if c.reader != nil {
		return c.reader.Close()
	}
	return nil
}

// GetEventHeader helper lấy header value từ message
func GetEventHeader(msg kafkago.Message, key string) string {
	for _, h := range msg.Headers {
		if h.Key == key {
			return string(h.Value)
		}
	}
	return ""
}

// UnmarshalEvent helper unmarshal message value vào struct
func UnmarshalEvent[T any](msg kafkago.Message, v *T) error {
	return json.Unmarshal(msg.Value, v)
}
