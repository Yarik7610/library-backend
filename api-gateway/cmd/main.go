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

	userMicroserviceHandler := core.ForwardTo(constants.USER_MICROSERVICE_SOCKET)
	catalogMicroserviceHandler := core.ForwardTo(constants.CATALOG_MICROSERVICE_SOCKET)

	userGroup := r.Group("")
	{
		userGroup.Use(middleware.AuthOptional())

		userGroup.POST(sharedconstants.SIGN_UP_ROUTE, userMicroserviceHandler)
		userGroup.POST(sharedconstants.SIGN_IN_ROUTE, userMicroserviceHandler)

		privateGroup := r.Group("")
		{
			privateGroup.Use(middleware.AuthRequired())

			privateGroup.GET(sharedconstants.ME_ROUTE, userMicroserviceHandler)
		}
	}

	catalogGroup := r.Group(sharedconstants.CATALOG_ROUTE)
	{
		catalogGroup.Use(middleware.AuthOptional())

		catalogGroup.GET(sharedconstants.CATEGORIES_ROUTE, catalogMicroserviceHandler)
		catalogGroup.GET(sharedconstants.CATEGORIES_ROUTE+"/:categoryName"+sharedconstants.BOOKS_ROUTE, catalogMicroserviceHandler)
		catalogGroup.GET(sharedconstants.AUTHORS_ROUTE+"/:authorID"+sharedconstants.BOOKS_ROUTE, catalogMicroserviceHandler)
		catalogGroup.GET(sharedconstants.BOOKS_ROUTE+sharedconstants.PREVIEW_ROUTE+"/:bookID", catalogMicroserviceHandler)
		catalogGroup.GET(sharedconstants.BOOKS_ROUTE+"/:bookID", catalogMicroserviceHandler)
		catalogGroup.GET(sharedconstants.BOOKS_ROUTE+sharedconstants.SEARCH_ROUTE, catalogMicroserviceHandler)
		catalogGroup.GET(sharedconstants.BOOKS_ROUTE+sharedconstants.NEW_ROUTE, catalogMicroserviceHandler)
		catalogGroup.GET(sharedconstants.BOOKS_ROUTE+sharedconstants.POPULAR_ROUTE, catalogMicroserviceHandler)
		catalogGroup.GET(sharedconstants.BOOKS_ROUTE+sharedconstants.VIEWS_ROUTE+"/:bookID", catalogMicroserviceHandler)

		adminGroup := catalogGroup.Group("")
		{
			adminGroup.Use(middleware.AuthRequired())

			adminGroup.DELETE(sharedconstants.BOOKS_ROUTE+"/:bookID", catalogMicroserviceHandler)
			adminGroup.POST(sharedconstants.BOOKS_ROUTE, catalogMicroserviceHandler)
			adminGroup.DELETE(sharedconstants.AUTHORS_ROUTE+"/:authorID", catalogMicroserviceHandler)
			adminGroup.POST(sharedconstants.AUTHORS_ROUTE, catalogMicroserviceHandler)
		}
	}

	subscriptionGroup := r.Group(sharedconstants.SUBSCRIPTIONS_ROUTE)
	{
		subscriptionGroup.Use(middleware.AuthRequired())

		subscriptionGroup.GET(sharedconstants.CATEGORIES_ROUTE, catalogMicroserviceHandler)
		subscriptionGroup.POST(sharedconstants.CATEGORIES_ROUTE, catalogMicroserviceHandler)
		subscriptionGroup.DELETE(sharedconstants.CATEGORIES_ROUTE+"/:categoryName", catalogMicroserviceHandler)
	}

	if err := r.Run(":" + config.Data.ServerPort); err != nil {
		zap.S().Fatalf("Start error on port %s: %v", config.Data.ServerPort, err)
	}
}
