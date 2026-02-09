package main

import (
	"github.com/Yarik7610/library-backend-common/microservice"
	"github.com/Yarik7610/library-backend/api-gateway/internal/app/router"
	"github.com/Yarik7610/library-backend/api-gateway/internal/core"
	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/config"
	"go.uber.org/zap"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
}

func main() {
	if err := config.Init(); err != nil {
		zap.S().Fatalf("Config load error: %v\n", err)
	}

	userServiceHandler := core.ForwardTo(microservice.USER_ADDRESS)
	catalogServiceHandler := core.ForwardTo(microservice.CATALOG_ADDRESS)
	subscriptionServiceHandler := core.ForwardTo(microservice.SUBSCRIPTIONS_ADDRESS)

	r := router.Register(userServiceHandler, catalogServiceHandler, subscriptionServiceHandler)

	if err := r.Run(":" + config.Data.ServerPort); err != nil {
		zap.S().Fatalf("Start error on port %s: %v", config.Data.ServerPort, err)
	}
}
