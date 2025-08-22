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

	userMicroserviceHandler := core.ForwardTo(constants.USER_MICROSERVICE_SOCKET)
	catalogMicroserviceHandler := core.ForwardTo(constants.CATALOG_MICROSERVICE_SOCKET)

	r.POST(sharedconstants.SIGN_UP_ROUTE, userMicroserviceHandler)
	r.POST(sharedconstants.SIGN_IN_ROUTE, userMicroserviceHandler)
	r.GET(sharedconstants.ME_ROUTE, userMicroserviceHandler)

	catalogRouter := r.Group(sharedconstants.CATALOG_ROUTE)
	{
		catalogRouter.GET(sharedconstants.CATEGORIES_ROUTE, catalogMicroserviceHandler)
		catalogRouter.GET(sharedconstants.CATEGORIES_ROUTE+"/:categoryName"+sharedconstants.BOOKS_ROUTE, catalogMicroserviceHandler)
		catalogRouter.GET(sharedconstants.AUTHORS_ROUTE+"/:authorID"+sharedconstants.BOOKS_ROUTE, catalogMicroserviceHandler)
		catalogRouter.GET(sharedconstants.BOOKS_ROUTE+sharedconstants.PREVIEW_ROUTE+"/:bookID", catalogMicroserviceHandler)
		catalogRouter.GET(sharedconstants.BOOKS_ROUTE+"/:bookID", catalogMicroserviceHandler)
		catalogRouter.GET(sharedconstants.BOOKS_ROUTE+sharedconstants.SEARCH_ROUTE, catalogMicroserviceHandler)
		catalogRouter.DELETE(sharedconstants.BOOKS_ROUTE+"/:bookID", catalogMicroserviceHandler)
		catalogRouter.POST(sharedconstants.BOOKS_ROUTE, catalogMicroserviceHandler)
		catalogRouter.DELETE(sharedconstants.AUTHORS_ROUTE+"/:authorID", catalogMicroserviceHandler)
		catalogRouter.POST(sharedconstants.AUTHORS_ROUTE, catalogMicroserviceHandler)
	}

	if err := r.Run(":" + config.Data.ServerPort); err != nil {
		zap.S().Fatalf("Start error on port %s: %v", config.Data.ServerPort, err)
	}
}
