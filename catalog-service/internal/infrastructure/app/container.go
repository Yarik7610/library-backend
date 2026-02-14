package app

import (
	"log"

	sharedKafka "github.com/Yarik7610/library-backend-common/broker/kafka"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/broker/kafka"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/observability/logging"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/storage/postgres"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/storage/redis"
)

type Container struct {
	Config         *config.Config
	Logger         *logging.Logger
	CatalogFeature *catalog.Feature
}

func NewContainer() *Container {
	config, err := config.Init()
	if err != nil {
		log.Fatalf("Config load error: %v\n", err)
	}

	logger := logging.NewLogger(config.Env)

	postgresDB, err := postgres.Connect(config)
	if err != nil {
		logger.Fatal("Postgres connect error", logging.Error(err))
	}
	logger.Info("Connected to Postgres and migrated it succesfully")

	redisClient, err := redis.Connect(config)
	if err != nil {
		logger.Fatal("Redis connect error", logging.Error(err))
	}
	logger.Info("Successfully connected to Redis")

	bookAddedWriter := kafka.NewWriter(sharedKafka.BOOK_ADDED_TOPIC)

	catalogFeature := catalog.NewFeature(logger, postgresDB, redisClient, bookAddedWriter)

	return &Container{
		Config:         config,
		Logger:         logger,
		CatalogFeature: catalogFeature,
	}
}
