package main

import (
	"github.com/Yarik7610/library-backend-common/sharedconstants"
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

	r.POST(sharedconstants.SIGN_UP_ROUTE, core.ForwardTo(constants.USER_MICROSERVICE_SOCKET, sharedconstants.SIGN_UP_ROUTE))
	r.POST(sharedconstants.SIGN_IN_ROUTE, core.ForwardTo(constants.USER_MICROSERVICE_SOCKET, sharedconstants.SIGN_IN_ROUTE))
	r.GET(sharedconstants.ME_ROUTE, core.ForwardTo(constants.USER_MICROSERVICE_SOCKET, sharedconstants.ME_ROUTE))

	catalogRouter := r.Group(sharedconstants.CATALOG_ROUTE)
	{
		catalogRouter.GET(sharedconstants.CATEGORIES_ROUTE, core.ForwardTo(constants.CATALOG_MICROSERVICE_SOCKET, sharedconstants.CATALOG_ROUTE+sharedconstants.CATEGORIES_ROUTE))
		catalogRouter.GET(sharedconstants.PREVIEW_ROUTE, core.ForwardTo(constants.CATALOG_MICROSERVICE_SOCKET, sharedconstants.CATALOG_ROUTE+sharedconstants.PREVIEW_ROUTE))
	}

	if err := r.Run(":" + config.Data.ServerPort); err != nil {
		zap.S().Fatalf("API-gateway start error on port %s: %v", config.Data.ServerPort, err)
	}
}
