package app

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/Yarik7610/library-backend/user-service/internal/feature/user"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/observability/logging"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/observability/tracing"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/storage/postgres"
)

type Container struct {
	Config          *config.Config
	Logger          *logging.Logger
	httpServer      *http.Server
	shutdownTracing func(context.Context) error
}

func NewContainer() *Container {
	config, err := config.Parse()
	if err != nil {
		log.Fatalf("Config parse error: %v\n", err)
	}

	logger := logging.NewLogger(config.Env)

	shutdownTracing, err := tracing.Init(config)
	if err != nil {
		logger.Fatal(context.Background(), "Tracing init error", logging.Error(err))
	}

	postgresDB, err := postgres.Connect(config)
	if err != nil {
		logger.Fatal(context.Background(), "Postgres connect error", logging.Error(err))
	}

	userFeature, err := user.NewFeature(config, logger, postgresDB)
	if err != nil {
		logger.Fatal(context.Background(), "User feature init error", logging.Error(err))
	}

	httpServer := &http.Server{
		Addr:    ":" + config.HTTPServerPort,
		Handler: userFeature.HTTPRouter,
	}

	return &Container{
		Config:          config,
		Logger:          logger,
		httpServer:      httpServer,
		shutdownTracing: shutdownTracing,
	}
}

func (c *Container) Start() error {
	err := c.httpServer.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func (c *Container) Stop(ctx context.Context) error {
	if err := c.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	if err := c.shutdownTracing(ctx); err != nil {
		return err
	}

	return nil
}
