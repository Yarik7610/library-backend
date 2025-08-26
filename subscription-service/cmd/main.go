package main

import (
	"github.com/Yarik7610/library-backend-common/sharedconstants"
	"github.com/Yarik7610/library-backend/catalog-service/config"
	"github.com/Yarik7610/library-backend/catalog-service/connect"
	"github.com/Yarik7610/library-backend/catalog-service/internal/controller"
	"github.com/Yarik7610/library-backend/catalog-service/internal/repository"
	"github.com/Yarik7610/library-backend/catalog-service/internal/service"

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

	db := connect.DB()

	userCategoryRepo := repository.NewUserCategoryRepository(db)
	subscriptionService := service.NewSubscriptionService(userCategoryRepo)
	subscriptionController := controller.NewSubscriptionController(subscriptionService)

	r := gin.Default()
	subscriptionGroup := r.Group(sharedconstants.SUBSCRIPTIONS_ROUTE)
	{
		subscriptionGroup.GET(sharedconstants.CATEGORIES_ROUTE, subscriptionController.GetSubscribedCategories)
		subscriptionGroup.POST(sharedconstants.CATEGORIES_ROUTE, subscriptionController.SubscribeCategory)
		subscriptionGroup.DELETE(sharedconstants.CATEGORIES_ROUTE+"/:categoryName", subscriptionController.UnsubscribeCategory)

		nonAPIGatewayGroup := subscriptionGroup.Group("")
		{
			nonAPIGatewayGroup.GET(sharedconstants.CATEGORIES_ROUTE+"/:categoryName", subscriptionController.GetCategorySubscribersEmails)
		}
	}

	if err := r.Run(":" + config.Data.ServerPort); err != nil {
		zap.S().Fatalf("Start error on port %s: %v", config.Data.ServerPort, err)
	}
}
