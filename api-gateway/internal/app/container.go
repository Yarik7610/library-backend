package app

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/Yarik7610/library-backend-common/microservice"
	"github.com/Yarik7610/library-backend/api-gateway/internal/app/router"
	"github.com/Yarik7610/library-backend/api-gateway/internal/core"
	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/observability/logging"
	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/observability/metrics"
	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/observability/tracing"
	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/swagger"
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
		log.Fatalf("Config load error: %v\n", err)
	}

	logger := logging.NewLogger(config.Env)

	shutdownTracing, err := tracing.Init(config)
	if err != nil {
		logger.Fatal(context.Background(), "Tracing init error", logging.Error(err))
	}

	userMicroserviceHandler := core.ForwardTo(logger, microservice.USER_ADDRESS)
	catalogMicroserviceHandler := core.ForwardTo(logger, microservice.CATALOG_ADDRESS)
	subscriptionMicroserviceHandler := core.ForwardTo(logger, microservice.SUBSCRIPTIONS_ADDRESS)
	metricsHandler, err := metrics.Init()
	if err != nil {
		logger.Fatal(context.Background(), "Metrics init error", logging.Error(err))
	}
	swaggerHandler := swagger.NewHandler(config, logger)

	router := router.Register(
		logger, config,
		metricsHandler,
		swaggerHandler,
		userMicroserviceHandler,
		catalogMicroserviceHandler,
		subscriptionMicroserviceHandler,
	)

	httpServer := &http.Server{
		Addr:    ":" + config.HTTPServerPort,
		Handler: router,
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
