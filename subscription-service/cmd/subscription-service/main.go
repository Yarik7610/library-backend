package main

import (
	"github.com/Yarik7610/library-backend-common/transport/http/route"
	repository "github.com/Yarik7610/library-backend/subscription-service/internal/feauture/subscription/repository/postgres"
	"github.com/Yarik7610/library-backend/subscription-service/internal/feauture/subscription/service"
	controller "github.com/Yarik7610/library-backend/subscription-service/internal/feauture/subscription/transport/http"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/storage/postgres"

	docs "github.com/Yarik7610/library-backend/subscription-service/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
}

func main() {
	if err := config.Init(); err != nil {
		zap.S().Fatalf("Config load error: %v\n", err)
	}

	postgresDB := postgres.Connect()

	userCategoryRepo := repository.NewUserCategoryRepository(postgresDB)
	subscriptionService := service.NewSubscriptionService(userCategoryRepo)
	subscriptionController := controller.NewSubscriptionController(subscriptionService)

	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	subscriptionGroup := r.Group(route.SUBSCRIPTIONS)
	{
		bookCategoryGroup := subscriptionGroup.Group(route.BOOKS + route.CATEGORIES)
		{
			bookCategoryGroup.GET("", subscriptionController.GetSubscribedCategories)
			bookCategoryGroup.POST("", subscriptionController.SubscribeCategory)
			bookCategoryGroup.DELETE("/:categoryName", subscriptionController.UnsubscribeCategory)
		}

		nonAPIGatewayGroup := bookCategoryGroup.Group("")
		{
			nonAPIGatewayGroup.GET("/:categoryName", subscriptionController.GetCategorySubscribersEmails)
		}
	}

	if err := r.Run(":" + config.Data.ServerPort); err != nil {
		zap.S().Fatalf("Start error on port %s: %v", config.Data.ServerPort, err)
	}
}
