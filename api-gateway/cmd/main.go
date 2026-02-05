package main

import (
	"github.com/Yarik7610/library-backend-common/sharedconstants"
	"github.com/Yarik7610/library-backend/api-gateway/config"
	docs "github.com/Yarik7610/library-backend/api-gateway/docs"
	"github.com/Yarik7610/library-backend/api-gateway/internal/core"
	"github.com/Yarik7610/library-backend/api-gateway/internal/middleware"
	"github.com/Yarik7610/library-backend/api-gateway/internal/utils"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	userMicroserviceHandler := core.ForwardTo(sharedconstants.USER_MICROSERVICE_SOCKET)
	catalogMicroserviceHandler := core.ForwardTo(sharedconstants.CATALOG_MICROSERVICE_SOCKET)
	subscriptionMicroserviceHandler := core.ForwardTo(sharedconstants.SUBSCRIPTIONS_MICROSERVICE_SOCKET)

	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger-json/doc.json", func(ctx *gin.Context) {
		userDoc := utils.FetchSwaggerJSON(sharedconstants.USER_MICROSERVICE_SOCKET)
		catalogDoc := utils.FetchSwaggerJSON(sharedconstants.CATALOG_MICROSERVICE_SOCKET)
		subDoc := utils.FetchSwaggerJSON(sharedconstants.SUBSCRIPTIONS_MICROSERVICE_SOCKET)

		mergedDoc := utils.MergeSwaggerDocs(userDoc, catalogDoc, subDoc)
		ctx.JSON(200, mergedDoc)
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger-json/doc.json")))

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
		catalogGroup.GET(sharedconstants.BOOKS_ROUTE+sharedconstants.CATEGORIES_ROUTE, catalogMicroserviceHandler)
		catalogGroup.GET(sharedconstants.BOOKS_ROUTE+sharedconstants.CATEGORIES_ROUTE+"/:categoryName", catalogMicroserviceHandler)
		catalogGroup.GET(sharedconstants.AUTHORS_ROUTE+"/:authorID"+sharedconstants.BOOKS_ROUTE, catalogMicroserviceHandler)
		catalogGroup.GET(sharedconstants.BOOKS_ROUTE+sharedconstants.PREVIEW_ROUTE+"/:bookID", catalogMicroserviceHandler)
		catalogGroup.GET(sharedconstants.BOOKS_ROUTE+"/:bookID", catalogMicroserviceHandler)
		catalogGroup.GET(sharedconstants.BOOKS_ROUTE+sharedconstants.SEARCH_ROUTE, catalogMicroserviceHandler)
		catalogGroup.GET(sharedconstants.BOOKS_ROUTE+sharedconstants.NEW_ROUTE, catalogMicroserviceHandler)
		catalogGroup.GET(sharedconstants.BOOKS_ROUTE+sharedconstants.POPULAR_ROUTE, catalogMicroserviceHandler)
		catalogGroup.GET(sharedconstants.BOOKS_ROUTE+sharedconstants.VIEWS_ROUTE+"/:bookID", catalogMicroserviceHandler)

		adminGroup := catalogGroup.Group("")
		{
			adminGroup.DELETE(sharedconstants.BOOKS_ROUTE+"/:bookID", catalogMicroserviceHandler)
			adminGroup.POST(sharedconstants.BOOKS_ROUTE, catalogMicroserviceHandler)
			adminGroup.DELETE(sharedconstants.AUTHORS_ROUTE+"/:authorID", catalogMicroserviceHandler)
			adminGroup.POST(sharedconstants.AUTHORS_ROUTE, catalogMicroserviceHandler)
		}
	}

	subscriptionGroup := r.Group(sharedconstants.SUBSCRIPTIONS_ROUTE)
	{
		subscriptionGroup.Use(middleware.AuthRequired())

		subscriptionGroup.GET(sharedconstants.CATEGORIES_ROUTE, subscriptionMicroserviceHandler)
		subscriptionGroup.POST(sharedconstants.CATEGORIES_ROUTE, subscriptionMicroserviceHandler)
		subscriptionGroup.DELETE(sharedconstants.CATEGORIES_ROUTE+"/:categoryName", subscriptionMicroserviceHandler)
	}

	if err := r.Run(":" + config.Data.ServerPort); err != nil {
		zap.S().Fatalf("Start error on port %s: %v", config.Data.ServerPort, err)
	}
}
