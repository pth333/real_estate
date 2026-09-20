package initialize

import (
	"context"
	"real_estate_be/internal/cron"
	"real_estate_be/internal/global"
	"real_estate_be/internal/wire"
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

	// Start cron jobs của luồng đặt cọc (auto reject 24h, nhắc lịch, quá hạn báo cáo)
	startDepositScheduler(ctx)

	// Init routes
	app := InitRouter()
	app.Listen(":8000")
}

// startDepositScheduler khởi động scheduler đặt cọc trong 1 goroutine riêng.
func startDepositScheduler(ctx context.Context) {
	depositService, err := wire.InitializeDepositService()
	if err != nil {
		panic(err)
	}

	scheduler := cron.NewDepositScheduler(depositService, global.Config.Deposit.CronIntervalSeconds)
	go scheduler.Start(ctx)
}

func RunCrawler() {
	LoadConfig()
	InitMysql()

	// Migrate DB
	MigrateDb(global.DB)
}
