package app

import (
	"log"

	"github.com/Yarik7610/library-backend-common/microservice"
	"github.com/Yarik7610/library-backend/api-gateway/internal/app/router"
	"github.com/Yarik7610/library-backend/api-gateway/internal/core"
	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/observability/logging"
	"github.com/gin-gonic/gin"
)

type Container struct {
	Config *config.Config
	Logger *logging.Logger
	Router *gin.Engine
}

func NewContainer() *Container {
	config, err := config.Init()
	if err != nil {
		log.Fatalf("Config load error: %v\n", err)
	}

	logger := logging.NewLogger(config.Env)

	userServiceHandler := core.ForwardTo(logger, microservice.USER_ADDRESS)
	catalogServiceHandler := core.ForwardTo(logger, microservice.CATALOG_ADDRESS)
	subscriptionServiceHandler := core.ForwardTo(logger, microservice.SUBSCRIPTIONS_ADDRESS)

	router := router.Register(
		logger, config,
		userServiceHandler,
		catalogServiceHandler,
		subscriptionServiceHandler,
	)

	return &Container{
		Config: config,
		Logger: logger,
		Router: router,
	}
}
