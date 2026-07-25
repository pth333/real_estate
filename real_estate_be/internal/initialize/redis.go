package initialize

import (
	"context"
	"fmt"
	"real_estate_be/internal/global"

	"github.com/redis/go-redis/v9"
)

func InitRedis() {
	cfg := global.Config.Redis
	if cfg.Addr == "" {
		cfg.Addr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
		DB:   cfg.DB,
	})

	// Ping để kiểm tra kết nối
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(fmt.Sprintf("Redis connection failed: %v", err))
	}

	global.RedisClient = rdb
}
