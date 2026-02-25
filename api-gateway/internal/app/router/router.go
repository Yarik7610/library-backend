package router

import (
	"net/http"

	"github.com/Yarik7610/library-backend-common/transport/http/route"
	"github.com/Yarik7610/library-backend/api-gateway/internal/app/middleware"
	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/observability/logging"
	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/swagger"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func Register(
	logger *logging.Logger,
	config *config.Config,
	metricsHandler http.Handler,
	swaggerHandler swagger.Handler,
	userMicroserviceHandler gin.HandlerFunc,
	catalogMicroserviceHandler gin.HandlerFunc,
	subscriptionMicroserviceHandler gin.HandlerFunc,
) *gin.Engine {
	r := gin.Default()

	r.Use(otelgin.Middleware(config.ServiceName,
		otelgin.WithGinFilter(func(c *gin.Context) bool {
			path := c.FullPath()
			return path != route.METRICS &&
				path != "/swagger-json/doc.json" &&
				path != "/swagger/*any"
		}),
	))
	r.Use(middleware.AuthContext(config))

	r.GET(route.METRICS, gin.WrapH(metricsHandler))

	swagger.RegisterRoutes(r, swaggerHandler)

	registerUserRoutes(r, userMicroserviceHandler)
	registerCatalogRoutes(r, catalogMicroserviceHandler)
	registerSubscriptionRoutes(r, subscriptionMicroserviceHandler)

	return r
}
