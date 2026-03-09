package cache

import (
	"context"
	"log"

	"github.com/example/learn-go-gin/internal/config"
	"github.com/redis/go-redis/v9"
)

// RedisClient 全域 Redis client 實例
var RedisClient *redis.Client

// Ctx 快取操作使用的 context
var Ctx = context.Background()

// InitRedis 初始化 Redis 連線
func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.GetRedisAddr(),
		Password: "", // 無密碼
		DB:       0,  // 使用預設 DB
	})

	// 測試連線
	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connected successfully")
}
