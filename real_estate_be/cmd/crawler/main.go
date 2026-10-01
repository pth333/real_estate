package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"real_estate_be/internal/crawler"
	provider "real_estate_be/internal/crawler/provider"
	"real_estate_be/internal/global"
	"real_estate_be/internal/initialize"
	kafkapkg "real_estate_be/pkg/kafka"
	"real_estate_be/pkg/vntime"
)

func main() {
	// Ép cả tiến trình dùng giờ VN (UTC+7) trước khi mở DB (VPS thường chạy UTC)
	vntime.SetDefault()

	// Load config + DB
	initialize.RunCrawler()

	// Init Kafka producer
	producer := kafkapkg.NewProducer()
	defer producer.Close()

	// Init crawler
	c := provider.NewBatDongSanCrawler()

	// Create scheduler chạy mỗi 30 phút
	scheduler := crawler.NewScheduler(30*time.Minute, c, global.DB, producer)

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("🏁 Crawler scheduler started. Press Ctrl+C to stop.")
	scheduler.Start(ctx)
	log.Println("🏁 Crawler stopped.")
}
