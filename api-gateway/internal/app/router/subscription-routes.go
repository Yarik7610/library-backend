package router

import (
	"github.com/Yarik7610/library-backend-common/transport/http/route"
	"github.com/Yarik7610/library-backend/api-gateway/internal/app/middleware"
	"github.com/Yarik7610/library-backend/api-gateway/internal/core"
	"github.com/gin-gonic/gin"
)

func registerSubscriptionRoutes(r *gin.Engine, subscriptionMicroserviceHandler gin.HandlerFunc) {
	subscriptionGroup := r.Group(route.SUBSCRIPTIONS)
	{
		bookCategoryGroup := subscriptionGroup.Group(route.BOOKS + route.CATEGORIES)
		bookCategoryGroup.Use(middleware.AuthRequired(), core.InjectHeaders())
		{
			bookCategoryGroup.GET("", subscriptionMicroserviceHandler)
			bookCategoryGroup.POST("", subscriptionMicroserviceHandler)
			bookCategoryGroup.DELETE("/:categoryName", subscriptionMicroserviceHandler)
		}
	}
}
