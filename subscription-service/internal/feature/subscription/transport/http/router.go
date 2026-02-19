package http

import (
	"net/http"

	"github.com/Yarik7610/library-backend-common/transport/http/route"
	"github.com/Yarik7610/library-backend/subscription-service/docs"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/config"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func NewRouter(config *config.Config, metricsHandler http.Handler, subscriptionHandler SubscriptionHandler) *gin.Engine {
	r := gin.Default()

	r.Use(otelgin.Middleware(config.ServiceName))
	r.GET("/metrics", gin.WrapH(metricsHandler))

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	subscriptionGroup := r.Group(route.SUBSCRIPTIONS)
	{
		bookCategoryGroup := subscriptionGroup.Group(route.BOOKS + route.CATEGORIES)
		{
			bookCategoryGroup.GET("", subscriptionHandler.GetUserSubscribedBookCategories)
			bookCategoryGroup.POST("", subscriptionHandler.SubscribeToBookCategory)
			bookCategoryGroup.DELETE("/:categoryName", subscriptionHandler.UnsubscribeFromBookCategory)
		}

		nonAPIGatewayGroup := bookCategoryGroup.Group("")
		{
			nonAPIGatewayGroup.GET("/:categoryName", subscriptionHandler.GetBookCategorySubscribedUserEmails)
		}
	}

	return r
}
