package initialize

import (
	"context"

	"real_estate_be/internal/global"
)

func Run() {
	LoadConfig()
	InitMysql()

	// Migrate DB
	MigrateDb(global.DB)

	// Init Kafka + SSE
	InitKafka()

	// Start Kafka consumers (background)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	StartKafkaConsumers(ctx, global.DB)

	// Init routes
	app := InitRouter()
	app.Listen(":8000")
}

func RunCrawler() {
	LoadConfig()
	InitMysql()

	// Migrate DB
	MigrateDb(global.DB)

	// Crawler được start từ cmd/crawler/main.go
}
