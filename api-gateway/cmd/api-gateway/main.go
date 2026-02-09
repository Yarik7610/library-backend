package main

import (
	"github.com/Yarik7610/library-backend-common/microservice"
	"github.com/Yarik7610/library-backend-common/transport/http/route"
	"github.com/Yarik7610/library-backend/api-gateway/internal/core"
	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/swagger"
	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/transport/http/middleware"
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

	userMicroserviceHandler := core.ForwardTo(microservice.USER_ADDRESS)
	catalogMicroserviceHandler := core.ForwardTo(microservice.CATALOG_ADDRESS)
	subscriptionMicroserviceHandler := core.ForwardTo(microservice.SUBSCRIPTIONS_ADDRESS)

	r := gin.Default()
	r.Use(middleware.AuthContext())

	swagger.InitRoutes(r)

	userGroup := r.Group("")
	{
		userGroup.POST(route.SIGN_UP, userMicroserviceHandler)
		userGroup.POST(route.SIGN_IN, userMicroserviceHandler)

		privateGroup := userGroup.Group("")
		privateGroup.Use(middleware.AuthRequired(), middleware.InjectHeaders())
		{
			privateGroup.GET(route.ME, userMicroserviceHandler)
		}
	}

	catalogGroup := r.Group(route.CATALOG)
	{
		catalogGroup.GET(route.BOOKS+route.CATEGORIES, catalogMicroserviceHandler)
		catalogGroup.GET(route.BOOKS+route.CATEGORIES+"/:categoryName", catalogMicroserviceHandler)
		catalogGroup.GET(route.AUTHORS+"/:authorID"+route.BOOKS, catalogMicroserviceHandler)
		catalogGroup.GET(route.BOOKS+"/:bookID"+route.PREVIEW, middleware.InjectHeaders(), catalogMicroserviceHandler)
		catalogGroup.GET(route.BOOKS+"/:bookID", catalogMicroserviceHandler)
		catalogGroup.GET(route.BOOKS+route.SEARCH, catalogMicroserviceHandler)
		catalogGroup.GET(route.BOOKS+route.NEW, catalogMicroserviceHandler)
		catalogGroup.GET(route.BOOKS+route.POPULAR, catalogMicroserviceHandler)
		catalogGroup.GET(route.BOOKS+"/:bookID"+route.VIEWS, catalogMicroserviceHandler)

		adminGroup := catalogGroup.Group("")
		adminGroup.Use(middleware.AuthRequired(), middleware.AdminRequired(), middleware.InjectHeaders())
		{
			adminGroup.DELETE(route.BOOKS+"/:bookID", catalogMicroserviceHandler)
			adminGroup.POST(route.BOOKS, catalogMicroserviceHandler)
			adminGroup.DELETE(route.AUTHORS+"/:authorID", catalogMicroserviceHandler)
			adminGroup.POST(route.AUTHORS, catalogMicroserviceHandler)
		}
	}

	subscriptionGroup := r.Group(route.SUBSCRIPTIONS)
	subscriptionGroup.Use(middleware.AuthRequired(), middleware.InjectHeaders())
	{
		subscriptionGroup.GET(route.CATEGORIES, subscriptionMicroserviceHandler)
		subscriptionGroup.POST(route.CATEGORIES, subscriptionMicroserviceHandler)
		subscriptionGroup.DELETE(route.CATEGORIES+"/:categoryName", subscriptionMicroserviceHandler)
	}

	if err := r.Run(":" + config.Data.ServerPort); err != nil {
		zap.S().Fatalf("Start error on port %s: %v", config.Data.ServerPort, err)
	}
}
