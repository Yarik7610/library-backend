package http

import (
	"github.com/Yarik7610/library-backend-common/transport/http/route"
	"github.com/Yarik7610/library-backend/subscription-service/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(subscriptionHandler SubscriptionHandler) *gin.Engine {
	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	subscriptionGroup := r.Group(route.SUBSCRIPTIONS)
	{
		bookCategoryGroup := subscriptionGroup.Group(route.BOOKS + route.CATEGORIES)
		{
			bookCategoryGroup.GET("", subscriptionHandler.GetUserBookCategories)
			bookCategoryGroup.POST("", subscriptionHandler.Create)
			bookCategoryGroup.DELETE("/:categoryName", subscriptionHandler.Delete)
		}

		nonAPIGatewayGroup := bookCategoryGroup.Group("")
		{
			nonAPIGatewayGroup.GET("/:categoryName", subscriptionHandler.GetCategorySubscribersEmails)
		}
	}

	return r
}
