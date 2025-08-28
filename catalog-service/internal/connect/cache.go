package connect

import (
	"context"
	"fmt"

	"github.com/Yarik7610/library-backend/catalog-service/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func Cache() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Data.RedisHost, config.Data.RedisPort),
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		zap.S().Fatalf("Redis connect error: %v", err)
	}
	zap.S().Info("Successfully connected to Redis")

	return rdb
}
