package router

import (
	"github.com/Yarik7610/library-backend-common/transport/http/route"
	"github.com/Yarik7610/library-backend/api-gateway/internal/app/middleware"
	"github.com/Yarik7610/library-backend/api-gateway/internal/core"
	"github.com/gin-gonic/gin"
)

func registerSubscriptionRoutes(r *gin.Engine, subscriptionServiceHandler gin.HandlerFunc) {
	subscriptionGroup := r.Group(route.SUBSCRIPTIONS)
	subscriptionGroup.Use(middleware.AuthRequired(), core.InjectHeaders())
	{
		subscriptionGroup.GET(route.CATEGORIES, subscriptionServiceHandler)
		subscriptionGroup.POST(route.CATEGORIES, subscriptionServiceHandler)
		subscriptionGroup.DELETE(route.CATEGORIES+"/:categoryName", subscriptionServiceHandler)
	}
}
