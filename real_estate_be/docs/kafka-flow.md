# Kafka Pipeline — kiến trúc & luồng xử lý

## Tổng quan

Dự án dùng Kafka làm event bus để xử lý dữ liệu bất động sản theo pipeline 3 tầng:

```
[Crawler] ──publish──▶ Topic crawled ──▶ [EnrichConsumer] ──publish──▶ Topic enriched ──▶ [NotifyConsumer (stub)]
```

Mỗi tầng là một consumer độc lập, giao tiếp qua Kafka topic. Crawler là producer gốc, các consumer vừa consume vừa produce cho tầng tiếp theo.

---

## Cấu hình Kafka

File: `config/config_local.yaml`

```yaml
kafka:
  brokers:
    - localhost:9092
  client_id: real-estate-be
  group_prefix: real-estate-be
  topics:
    real_estate_crawled: real_estate.crawled.v1
    real_estate_enriched: real_estate.enriched.v1
    real_estate_notified: real_estate.notified.v1
```

Cấu hình được nạp vào `global.Config.Kafka` (`internal/global/global.go:27-38`).

### Các topic Kafka

| Config key | Topic name | Event type | Mô tả |
|---|---|---|---|
| `real_estate_crawled` | `real_estate.crawled.v1` | `RealEstateCrawledEvent` | Raw data sau khi crawl |
| `real_estate_enriched` | `real_estate.enriched.v1` | `RealEstateEnrichedEvent` | Data đã enrich (phân loại, toạ độ) |
| `real_estate_notified` | `real_estate.notified.v1` | `RealEstateNotifiedEvent` | Log notification (chưa dùng) |

---

## Event types

File: `pkg/kafka/event.go`

3 struct event, tất cả đều embed `BaseEvent`:

### BaseEvent — fields chung

```go
type BaseEvent struct {
    EventType string    `json:"event_type"`
    Source    string    `json:"source"`
    Version   string    `json:"version"`
    Timestamp time.Time `json:"timestamp"`
    TraceID   string    `json:"trace_id,omitempty"`
}
```

### RealEstateCrawledEvent

| Field | Kiểu | Ghi chú |
|---|---|---|
| SourceURL | string | Key Kafka message |
| Title, Address, District, City | string | |
| PriceVND | float64 | |
| Acreage | float64 | |
| PricePerM2 | float64 | |
| CrawledAt | time.Time | |
| PublishedAt | *time.Time | |

### RealEstateEnrichedEvent

Thừa kế `RealEstateCrawledEvent`, thêm:

| Field | Kiểu |
|---|---|
| TypeOfRealEstate | string |
| Latitude | *float64 |
| Longitude | *float64 |

### RealEstateNotifiedEvent

| Field | Kiểu |
|---|---|
| SourceURL | string |
| Channel | string (email, sms, webhook) |
| Recipients | int |
| Success | bool |

Mỗi event implement các interface: `GetKey()`, `GetEventType()`, `GetSource()`, `GetVersion()` — dùng để tự động gắn Kafka headers khi publish.

### Event headers

```go
HeaderEventType   = "x-event-type"
HeaderSource      = "x-source"
HeaderVersion     = "x-event-version"
HeaderTimestamp   = "x-timestamp"
HeaderTraceID     = "x-trace-id"
HeaderContentType = "content-type"
```

---

## Producer

File: `pkg/kafka/producer.go`

Dùng thư viện [`segmentio/kafka-go`](https://github.com/segmentio/kafka-go).

```go
type Producer struct {
    writer *kafkago.Writer
    mu     sync.Mutex
    closed bool
}
```

### Khởi tạo: `NewProducer()`

- Đọc brokers + default topic từ `global.Config`
- Tạo `Writer` với:
  - `AllowAutoTopicCreation: true` — tự tạo topic nếu chưa có
  - `Balancer: &LeastBytes{}` — partition theo dung lượng nhỏ nhất
  - `RequiredAcks: RequireOne` — chờ leader ack
  - `BatchTimeout: 10ms`, `BatchSize: 50`

### Publish message

```go
func (p *Producer) Publish(ctx context.Context, key string, event any) error
```

1. Check closed state (thread-safe qua mutex)
2. `json.Marshal(event)` → []byte
3. `buildMessage(key, data, event)` → gắn headers (content-type, event-type, source, version)
4. `writer.WriteMessages(ctx, msg)`

### `buildMessage()` — gắn headers tự động

Dùng type assertion để check interface `GetEventType()`, `GetSource()`, `GetVersion()` — nếu event implement thì thêm header tương ứng.

### `PublishBatch(ctx, keys, events)`

Gửi nhiều message cùng lúc (atomic write).

### `SetTopic(topic)`

Cho phép 1 Producer instance đổi topic — dùng trong `EnrichConsumer` khi publish sang enriched topic.

### Thread safety

- `sync.Mutex` bảo vệ `closed` flag + writer
- `Close()` idempotent

---

## Consumer

File: `pkg/kafka/consumer.go`

Base consumer reusable (dùng cho cả Enrich và Notify).

### `ConsumerConfig`

| Field | Mô tả |
|---|---|
| Topic | Tên topic cần consume |
| GroupSuffix | Suffix cho consumer group ID |
| Handler | Callback `func(ctx, msg) error` |
| Concurrency | Số goroutine xử lý song song (default 1) |

### Khởi tạo: `NewConsumer(cfg)`

- Group ID = `group_prefix` + `-` + `group_suffix` (VD: `real-estate-be-enrich`)
- `StartOffset: FirstOffset` — đọc từ đầu partition nếu chưa có offset
- `CommitInterval: 1s` — auto-commit

### Chế độ xử lý

| Chế độ | Khi nào |
|---|---|
| `loop()` — tuần tự | `concurrency = 1` |
| `loopConcurrent()` — song song | `concurrency > 1` (dùng semaphore) |

Cả 2 chế độ đều:
1. `FetchMessage(ctx)` — lấy message
2. Gọi `handler(ctx, msg)`
3. `CommitMessages(ctx, msg)` — commit offset

**Không retry** — handler lỗi vẫn commit để không block consumer.

### Helpers

```go
func GetEventHeader(msg, key) string    // lấy header value
func UnmarshalEvent[T](msg, *T) error   // generic unmarshal
```

---

## EnrichConsumer

File: `internal/kafka/enrich.go`

Consumer đầu tiên trong pipeline.

```
Topic "real_estate.crawled.v1"
    │
    ▼
EnrichConsumer (concurrency = 2)
    │
    ├──► Parse event
    ├──► classifyType() → phân loại BĐS
    ├──► Lưu DB (real_estate_enriched) — UPSERT theo source_url
    │
    └──► Publish RealEstateEnrichedEvent → Topic "real_estate.enriched.v1"
```

### Luồng xử lý

1. Check header `x-event-type == "real_estate.crawled.v1"`
2. Unmarshal → `RealEstateCrawledEvent`
3. `classifyType(title, price)` — hiện là stub, luôn trả `"apartment"`
4. Lưu `model.RealEstateEnriched` vào DB với `clause.OnConflict` (UPSERT)
5. Build `RealEstateEnrichedEvent` (copy fields + enriched fields)
6. `producer.Publish(ctx, SourceURL, enrichedEvent)` — publish lên enriched topic
7. Nếu lỗi → log + return nil (không block consumer)

---

## NotifyConsumer

File: `internal/kafka/notify.go`

Consumer thứ hai — **hiện là stub, chưa implement thật**.

```
Topic "real_estate.enriched.v1"
    │
    ▼
NotifyConsumer (concurrency = 1)
    │
    ├──► Parse RealEstateEnrichedEvent
    ├──► Log ra console
    │
    └──► TODO: gửi email/SMS/webhook
         TODO: publish RealEstateNotifiedEvent
```

---

## Khởi động consumers

File: `internal/initialize/kafka.go`

```go
func StartKafkaConsumers(ctx context.Context)
```

- Kiểm tra `len(brokers) == 0` → skip nếu không có Kafka
- Start EnrichConsumer trong 1 goroutine
- Start NotifyConsumer trong 1 goroutine
- Được gọi từ `main()` hoặc `initialize` module

---

## Crawler → Kafka

File: `internal/crawler/run.go` — hàm `handlePageResult()`

```go
// 1. Crawl page → raw data
// 2. Lưu vào DB (real_estate table)
// 3. Publish event lên Kafka:
event := kafka.NewRealEstateCrawledEvent(item)
producer.Publish(ctx, item.SourceURL, event)
```

- Key message = `SourceURL` → cùng URL luôn vào cùng partition (ordering)
- Value = JSON của `RealEstateCrawledEvent`

---

## Luồng đầy đủ

```
Crawler (run.go)
  │
  ├──► DB (lưu raw real_estate)
  │
  └──► Kafka Producer ──► Topic "real_estate.crawled.v1"
                                │
                           EnrichConsumer
                                │
                    ┌───────────┼───────────┐
                    │           │           │
                  Parse     classifyType  Lưu DB
                  event     (stub)        enriched
                    │           │           │
                    └───────────┼───────────┘
                                │
                    Kafka Producer ──► Topic "real_estate.enriched.v1"
                                            │
                                       NotifyConsumer (STUB)
                                            │
                                   TODO: send email/SMS
                                            │
                                   TODO: publish notified event
```

---

## Điểm đáng chú ý

| Đặc điểm | Giải thích |
|---|---|
| **Không retry** | Producer/Consumer đều không retry. Lỗi → log + bỏ qua + commit tiếp |
| **UPSERT** | EnrichConsumer dùng `OnConflict` — không insert trùng `source_url` |
| **classifyType là stub** | Luôn trả `"apartment"` |
| **NotifyConsumer chưa làm** | Mới chỉ log |
| **Headers tự động** | Dùng interface `GetEventType()`, `GetSource()`, `GetVersion()` để gắn headers — không cần code riêng từng event |
| **Event type có version** | `real_estate.crawled.v1`, `enriched.v1`, `notified.v1` |
| **Key = SourceURL** | Đảm bảo ordering cho cùng 1 bài đăng |
