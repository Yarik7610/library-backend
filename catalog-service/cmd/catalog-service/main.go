package main

import (
	"github.com/Yarik7610/library-backend-common/broker"
	"github.com/Yarik7610/library-backend-common/sharedconstants"

	"github.com/Yarik7610/library-backend/catalog-service/internal/feature"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/storage/postgres"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/storage/redis"
	"go.uber.org/zap"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
}

func main() {
	err := config.Init()
	if err != nil {
		zap.S().Fatalf("Config load error: %v\n", err)
	}

	postgresDB := postgres.Connect()
	redisClient := redis.Connect()
	bookAddedWriter := broker.NewWriter(sharedconstants.BOOK_ADDED_TOPIC)

	catalogFeature := feature.NewCatalog(postgresDB, redisClient, bookAddedWriter)

	zap.S().Infof("Starting catalog-service on port: '%s'", config.Data.ServerPort)

	if err := catalogFeature.HTTPRouter.Run(":" + config.Data.ServerPort); err != nil {
		zap.S().Fatalf("Start error on port %s: %v", config.Data.ServerPort, err)
	}
}
