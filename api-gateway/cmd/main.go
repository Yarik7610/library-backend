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

	r.POST(sharedconstants.SIGN_UP_ROUTE, core.ForwardTo(constants.USER_MICROSERVICE_SOCKET))
	r.POST(sharedconstants.SIGN_IN_ROUTE, core.ForwardTo(constants.USER_MICROSERVICE_SOCKET))
	r.GET(sharedconstants.ME_ROUTE, core.ForwardTo(constants.USER_MICROSERVICE_SOCKET))

	catalogRouter := r.Group(sharedconstants.CATALOG_ROUTE)
	{
		catalogRouter.GET(sharedconstants.CATEGORIES_ROUTE, core.ForwardTo(constants.CATALOG_MICROSERVICE_SOCKET))
		catalogRouter.GET(sharedconstants.PREVIEW_ROUTE+"/:bookID", core.ForwardTo(constants.CATALOG_MICROSERVICE_SOCKET))
		catalogRouter.GET(sharedconstants.BOOKS_ROUTE+"/:authorName", core.ForwardTo(constants.CATALOG_MICROSERVICE_SOCKET))
		catalogRouter.GET(sharedconstants.BOOKS_ROUTE+sharedconstants.SEARCH_ROUTE, core.ForwardTo(constants.CATALOG_MICROSERVICE_SOCKET))
	}

	if err := r.Run(":" + config.Data.ServerPort); err != nil {
		zap.S().Fatalf("API-gateway start error on port %s: %v", config.Data.ServerPort, err)
	}
}
