package redis

import (
	"context"
	"fmt"

	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/config"
	"github.com/redis/go-redis/v9"
)

func Connect(config *config.Config) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return redisClient, nil
}
