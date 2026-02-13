package app

import (
	"log"

	"github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/observability/logging"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/storage/postgres"
)

type Container struct {
	Config              *config.Config
	Logger              *logging.Logger
	SubscriptionFeature *subscription.Feature
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

	subscriptionFeature := subscription.NewFeature(logger, postgresDB)

	return &Container{
		Config:              config,
		Logger:              logger,
		SubscriptionFeature: subscriptionFeature,
	}
}
