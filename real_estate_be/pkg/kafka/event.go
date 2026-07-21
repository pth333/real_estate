package kafka

import (
	"time"

	model "real_estate_be/internal/models"
)

// Header keys.
const (
	HeaderEventType   = "x-event-type"
	HeaderSource      = "x-source"
	HeaderVersion     = "x-event-version"
	HeaderTimestamp   = "x-timestamp"
	HeaderTraceID     = "x-trace-id"
	HeaderContentType = "content-type"
)

// ── Interfaces cho buildMessage tự gắn headers ──

type EventTyper interface{ GetEventType() string }
type EventSourcer interface{ GetSource() string }
type EventVersion interface{ GetVersion() string }
type EventKeyer interface{ GetKey() string }

// ── BaseEvent ──

type BaseEvent struct {
	EventType string    `json:"event_type"`
	Source    string    `json:"source"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
	TraceID   string    `json:"trace_id,omitempty"`
}

func (e BaseEvent) GetEventType() string { return e.EventType }
func (e BaseEvent) GetSource() string    { return e.Source }
func (e BaseEvent) GetVersion() string   { return e.Version }

// ── Event type constants ──

const (
	EventTypeCrawled  = "real_estate.crawled.v1"
	EventTypeEnriched = "real_estate.enriched.v1"
	EventTypeNotified = "real_estate.notified.v1"

	SourceCrawler = "crawler"
	SourceEnrich  = "enrich-consumer"
	SourceNotify  = "notify-consumer"

	VersionV1 = "v1"
)

// ── RealEstateCrawledEvent ──

type RealEstateCrawledEvent struct {
	BaseEvent

	SourceURL   string     `json:"source_url"`
	Title       string     `json:"title"`
	Address     string     `json:"address"`
	District    string     `json:"district"`
	City        string     `json:"city"`
	PriceVND    float64    `json:"price_vnd"`
	Acreage     float64    `json:"acreage"`
	PricePerM2  float64    `json:"price_per_m2"`
	CrawledAt   time.Time  `json:"crawled_at"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
}

func (e RealEstateCrawledEvent) GetKey() string { return e.SourceURL }

func NewRealEstateCrawledEvent(m model.RealEstate) RealEstateCrawledEvent {
	return RealEstateCrawledEvent{
		BaseEvent: BaseEvent{
			EventType: EventTypeCrawled,
			Source:    SourceCrawler,
			Version:   VersionV1,
			Timestamp: time.Now(),
		},
		SourceURL:   m.SourceURL,
		Title:       m.Title,
		Address:     m.Address,
		District:    m.District,
		City:        m.City,
		PriceVND:    m.PriceVND,
		Acreage:     m.Acreage,
		PricePerM2:  m.PricePerM2,
		CrawledAt:   m.CrawledAt,
		PublishedAt: nil,
	}
}

// ── RealEstateEnrichedEvent ──

type RealEstateEnrichedEvent struct {
	BaseEvent

	SourceURL   string     `json:"source_url"`
	Title       string     `json:"title"`
	Address     string     `json:"address"`
	District    string     `json:"district"`
	City        string     `json:"city"`
	PriceVND    float64    `json:"price_vnd"`
	Acreage     float64    `json:"acreage"`
	PricePerM2  float64    `json:"price_per_m2"`
	CrawledAt   time.Time  `json:"crawled_at"`
	PublishedAt *time.Time `json:"published_at,omitempty"`

	// Enriched fields
	TypeOfRealEstate string   `json:"type_of_real_estate"`
	Latitude         *float64 `json:"latitude,omitempty"`
	Longitude        *float64 `json:"longitude,omitempty"`
}

func (e RealEstateEnrichedEvent) GetKey() string { return e.SourceURL }

func NewRealEstateEnrichedEvent(crawled RealEstateCrawledEvent, typeStr string) RealEstateEnrichedEvent {
	return RealEstateEnrichedEvent{
		BaseEvent: BaseEvent{
			EventType: EventTypeEnriched,
			Source:    SourceEnrich,
			Version:   VersionV1,
			Timestamp: time.Now(),
		},
		SourceURL:   crawled.SourceURL,
		Title:       crawled.Title,
		Address:     crawled.Address,
		District:    crawled.District,
		City:        crawled.City,
		PriceVND:    crawled.PriceVND,
		Acreage:     crawled.Acreage,
		PricePerM2:  crawled.PricePerM2,
		CrawledAt:   crawled.CrawledAt,
		PublishedAt: crawled.PublishedAt,

		TypeOfRealEstate: typeStr,
	}
}

// ── RealEstateNotifiedEvent ──

type RealEstateNotifiedEvent struct {
	BaseEvent

	SourceURL  string `json:"source_url"`
	Channel    string `json:"channel"`
	Recipients int    `json:"recipients"`
	Success    bool   `json:"success"`
}

func (e RealEstateNotifiedEvent) GetKey() string { return e.SourceURL }
