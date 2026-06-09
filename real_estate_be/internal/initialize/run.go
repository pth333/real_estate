package initialize

import (
	"context"
	"log"
	"os/signal"
	"syscall"
)

func Run() {
	LoadConfig()
	InitMysql()

	// Graceful shutdown context
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	// Start Kafka consumers (non-blocking)
	StartKafkaConsumers(ctx)

	// Start Fiber server
	app := InitRouter()

	// Shutdown server khi context done
	go func() {
		<-ctx.Done()
		log.Println("⏹ Shutting down server...")
		_ = app.Shutdown()
	}()

	if err := app.Listen(":8000"); err != nil {
		log.Printf("Server error: %v", err)
	}
}

func RunCrawler() {
	LoadConfig()
	InitMysql()
}
