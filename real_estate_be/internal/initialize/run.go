package initialize

import (
	"context"
	"real_estate_be/internal/global"
)

func Run() {
	LoadConfig()
	InitMysql()
	InitRedis()
	InitRecommendation()
	InitS3()

	// Migrate DB
	MigrateDb(global.DB)

	// Init Kafka + SSE
	InitKafka()

	// Start Kafka consumers (background)
	ctx := context.Background()
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
}
