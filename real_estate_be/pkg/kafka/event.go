package kafka

import (
	"time"

	model "real_estate_be/internal/models"
)

// ============================================================
// Event types — versioned semantic names used as message Key
// qua Kafka headers để consumer routing không phụ thuộc topic name
// ============================================================

const (
	// Crawled: emitted bởi crawler sau khi scrape được dữ liệu
	EventTypeCrawled = "real_estate.crawled.v1"

	// Enriched: emitted bởi enrichment worker sau khi enrich dữ liệu
	EventTypeEnriched = "real_estate.enriched.v1"

	// Notified: emitted sau khi gửi notification thành công
	EventTypeNotified = "real_estate.notified.v1"
)

// ============================================================
// Event headers — trace context cho toàn bộ pipeline
// ============================================================

const (
	HeaderEventType   = "x-event-type"
	HeaderSource      = "x-source"
	HeaderVersion     = "x-event-version"
	HeaderTimestamp   = "x-timestamp"
	HeaderTraceID     = "x-trace-id"
	HeaderContentType = "content-type"
)

const ContentTypeJSON = "application/json"

// ============================================================
// BaseEvent — các field chung cho mọi message
// ============================================================

type BaseEvent struct {
	EventType string    `json:"event_type"`
	Source    string    `json:"source"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
	TraceID   string    `json:"trace_id,omitempty"`
}

// ============================================================
// RealEstateCrawledEvent — từ crawler → raw data
// ============================================================

type RealEstateCrawledEvent struct {
	BaseEvent
	SourceURL   string    `json:"source_url"`
	Title       string    `json:"title"`
	Address     string    `json:"address"`
	District    string    `json:"district"`
	City        string    `json:"city"`
	PriceVND    float64   `json:"price_vnd"`
	Acreage     float64   `json:"acreage"`
	PricePerM2  float64   `json:"price_per_m2"`
	CrawledAt   time.Time `json:"crawled_at"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
}

// GetKey trả về Kafka message key (SourceURL)
func (e RealEstateCrawledEvent) GetKey() string { return e.SourceURL }

// GetEventType trả về event type
func (e RealEstateCrawledEvent) GetEventType() string { return e.EventType }

// GetSource trả về nguồn dữ liệu
func (e RealEstateCrawledEvent) GetSource() string { return e.Source }

// GetVersion trả về version
func (e RealEstateCrawledEvent) GetVersion() string { return e.Version }

// NewRealEstateCrawledEvent tạo event từ DB model sau khi đã lưu
func NewRealEstateCrawledEvent(item model.RealEstateModel) RealEstateCrawledEvent {
	return RealEstateCrawledEvent{
		BaseEvent: BaseEvent{
			EventType: EventTypeCrawled,
			Source:    item.Source,
			Version:   "1.0",
			Timestamp: time.Now(),
		},
		SourceURL:   item.SourceURL,
		Title:       item.Title,
		Address:     item.Address,
		District:    item.District,
		City:        item.City,
		PriceVND:    item.PriceVND,
		Acreage:     item.Acreage,
		PricePerM2:  item.PricePerM2,
		CrawledAt:   item.CrawledAt,
		PublishedAt: item.PublishedAt,
	}
}

// ============================================================
// RealEstateEnrichedEvent — sau enrichment step
// ============================================================

type RealEstateEnrichedEvent struct {
	BaseEvent
	SourceURL       string   `json:"source_url"`
	Title           string   `json:"title"`
	Address         string   `json:"address"`
	District        string   `json:"district"`
	City            string   `json:"city"`
	PriceVND        float64  `json:"price_vnd"`
	Acreage         float64  `json:"acreage"`
	PricePerM2      float64  `json:"price_per_m2"`
	TypeOfRealEstate string  `json:"type_of_real_estate"`
	Latitude        *float64 `json:"latitude,omitempty"`
	Longitude       *float64 `json:"longitude,omitempty"`
}

func (e RealEstateEnrichedEvent) GetKey() string      { return e.SourceURL }
func (e RealEstateEnrichedEvent) GetEventType() string  { return e.EventType }
func (e RealEstateEnrichedEvent) GetSource() string     { return e.Source }
func (e RealEstateEnrichedEvent) GetVersion() string    { return e.Version }

// ============================================================
// RealEstateNotifiedEvent — log notification đã gửi
// ============================================================

type RealEstateNotifiedEvent struct {
	BaseEvent
	SourceURL  string `json:"source_url"`
	Channel    string `json:"channel"` // email, sms, webhook, ...
	Recipients int    `json:"recipients"`
	Success    bool   `json:"success"`
}

func (e RealEstateNotifiedEvent) GetKey() string        { return e.SourceURL }
func (e RealEstateNotifiedEvent) GetEventType() string  { return e.EventType }
func (e RealEstateNotifiedEvent) GetSource() string     { return e.Source }
func (e RealEstateNotifiedEvent) GetVersion() string    { return e.Version }
