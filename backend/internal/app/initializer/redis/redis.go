package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"takeaway-go/internal/app/config"
	"takeaway-go/pkg/logger"
)

var RedisClient *redis.Client

// InitRedis 初始化 Redis 连接
func InitRedis() {
	conf := config.AppConf.Redis
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       conf.DB,
	})

	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	logger.Log.Info("Redis connected successfully")
}
