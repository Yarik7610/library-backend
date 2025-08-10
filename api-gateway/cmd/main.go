package main

import (
	"github.com/Yarik7610/library-backend/api-gateway/config"
	"github.com/Yarik7610/library-backend/api-gateway/internal/constants"
	"github.com/Yarik7610/library-backend/api-gateway/internal/core"
	"github.com/Yarik7610/library-backend/api-gateway/internal/middleware"
	"github.com/gin-gonic/gin"
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

	r := gin.Default()
	r.Use(middleware.AuthMiddleware())

	r.POST(constants.SIGN_UP_ROUTE, core.ForwardTo(constants.USER_MICROSERVICE_SOCKET, constants.SIGN_UP_ROUTE))
	r.POST(constants.SIGN_IN_ROUTE, core.ForwardTo(constants.USER_MICROSERVICE_SOCKET, constants.SIGN_IN_ROUTE))
	r.GET(constants.ME_ROUTE, core.ForwardTo(constants.USER_MICROSERVICE_SOCKET, constants.ME_ROUTE))

	if err := r.Run(":" + config.Data.ServerPort); err != nil {
		zap.S().Fatalf("API-gateway start error on port %s: %v", config.Data.ServerPort, err)
	}
}
