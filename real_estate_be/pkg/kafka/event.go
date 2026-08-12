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
	EventTypeCrawled    = "real_estate.crawled.v1"
	EventTypeEnriched   = "real_estate.enriched.v1"
	EventTypeNewListing = "real_estate.new_listing.v1"
	EventTypeNotified   = "real_estate.notified.v1"

	SourceCrawler = "crawler"
	SourceEnrich  = "enrich-consumer"
	SourceApp     = "real-estate-app"
	SourceNotify  = "notify-consumer"

	VersionV1 = "v1"
)

// ── RealEstateCrawledEvent (Giữ lại cho crawler) ──

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

// ── RealEstateEnrichedEvent (Giữ lại cho enricher) ──

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

	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
}

func (e RealEstateEnrichedEvent) GetKey() string { return e.SourceURL }

// ── RealEstateNewListingEvent (Mới - Chuyên cho Notify) ──

type RealEstateNewListingEvent struct {
	BaseEvent

	ListingID uint64  `json:"listing_id"`
	Title     string  `json:"title"`
	Address   string  `json:"address"`
	PriceVND  float64 `json:"price_vnd"`
	Acreage   float64 `json:"acreage"`
	Slug      string  `json:"slug"`
}

func (e RealEstateNewListingEvent) GetKey() string {
	return time.Now().Format("20060102") // Partition key theo ngày hoặc ID
}

func NewRealEstateNewListingEvent(m model.RealEstate) RealEstateNewListingEvent {
	return RealEstateNewListingEvent{
		BaseEvent: BaseEvent{
			EventType: EventTypeNewListing,
			Source:    SourceApp,
			Version:   VersionV1,
			Timestamp: time.Now(),
		},
		ListingID: m.ID,
		Title:     m.Title,
		Address:   m.Address,
		PriceVND:  m.PriceVND,
		Acreage:   m.Acreage,
		Slug:      m.Slug,
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
