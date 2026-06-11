package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"real_estate_be/internal/global"

	kafkago "github.com/segmentio/kafka-go"
	"golang.org/x/sync/semaphore"
)

// HandlerFunc nhận message và trả về lỗi nếu xử lý thất bại.
type HandlerFunc func(ctx context.Context, msg kafkago.Message) error

// ConsumerConfig ...
type ConsumerConfig struct {
	Topic       string
	GroupSuffix string
	Handler     HandlerFunc
	Concurrency int // số goroutine xử lý song song (mặc định 1 -> tuần tự)
}

type Consumer struct {
	reader *kafkago.Reader
	cfg    ConsumerConfig
}

func NewConsumer(cfg ConsumerConfig) *Consumer {
	groupID := global.Config.Kafka.GroupPrefix + "-" + cfg.GroupSuffix
	if cfg.Concurrency < 1 {
		cfg.Concurrency = 1
	}

	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:        global.Config.Kafka.Brokers,
		GroupID:        groupID,
		Topic:          cfg.Topic,
		StartOffset:    kafkago.FirstOffset,
		CommitInterval: time.Second,
		MaxWait:        3 * time.Second,
	})

	return &Consumer{reader: reader, cfg: cfg}
}

// Start bắt đầu consume vô hạn, block goroutine.
func (c *Consumer) Start(ctx context.Context) {
	log.Printf("🔄 [Kafka] consumer starting: topic=%s group=%s concurrency=%d",
		c.cfg.Topic, c.reader.Config().GroupID, c.cfg.Concurrency)

	if c.cfg.Concurrency <= 1 {
		c.loop(ctx)
	} else {
		c.loopConcurrent(ctx)
	}
}

// loop — tuần tự, từng message một.
func (c *Consumer) loop(ctx context.Context) {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("❌ [Kafka] FetchMessage error: %v", err)
			continue
		}

		if err := c.cfg.Handler(ctx, msg); err != nil {
			log.Printf("⚠️ [Kafka] handler error (will commit anyway): %v", err)
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("❌ [Kafka] CommitMessages error: %v", err)
		}
	}
}

// loopConcurrent — song song, dùng semaphore giới hạn concurrency.
func (c *Consumer) loopConcurrent(ctx context.Context) {
	sem := semaphore.NewWeighted(int64(c.cfg.Concurrency))

	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("❌ [Kafka] FetchMessage error: %v", err)
			continue
		}

		if err := sem.Acquire(ctx, 1); err != nil {
			return
		}

		m := msg
		go func() {
			defer sem.Release(1)

			if err := c.cfg.Handler(ctx, m); err != nil {
				log.Printf("⚠️ [Kafka] handler error: %v", err)
			}

			if err := c.reader.CommitMessages(ctx, m); err != nil {
				log.Printf("❌ [Kafka] CommitMessages error: %v", err)
			}
		}()
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

// ── Helpers ──

func GetEventHeader(msg kafkago.Message, key string) string {
	for _, h := range msg.Headers {
		if h.Key == key {
			return string(h.Value)
		}
	}
	return ""
}

func UnmarshalEvent[T any](msg kafkago.Message, v *T) error {
	return json.Unmarshal(msg.Value, v)
}

func UnmarshalMsg[T any](msg kafkago.Message) (*T, error) {
	var v T
	if err := json.Unmarshal(msg.Value, &v); err != nil {
		return nil, fmt.Errorf("unmarshal %T: %w", v, err)
	}
	return &v, nil
}
