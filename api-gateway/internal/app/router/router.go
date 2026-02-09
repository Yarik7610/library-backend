package router

import (
	"github.com/Yarik7610/library-backend/api-gateway/internal/app/middleware"
	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/swagger"
	"github.com/gin-gonic/gin"
)

func Register(
	userServiceHandler gin.HandlerFunc,
	catalogServiceHandler gin.HandlerFunc,
	subscriptionServiceHandler gin.HandlerFunc,
) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.AuthContext())

	swagger.RegisterRoutes(r)

	registerUserRoutes(r, userServiceHandler)
	registerCatalogRoutes(r, catalogServiceHandler)
	registerSubscriptionRoutes(r, subscriptionServiceHandler)

	return r
}
