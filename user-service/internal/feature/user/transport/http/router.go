package http

import (
	"net/http"

	"github.com/Yarik7610/library-backend-common/transport/http/route"
	"github.com/Yarik7610/library-backend/user-service/docs"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/config"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func NewRouter(config *config.Config, metricsHandler http.Handler, userHandler UserHandler) *gin.Engine {
	r := gin.Default()

	r.Use(otelgin.Middleware(config.ServiceName,
		otelgin.WithGinFilter(func(c *gin.Context) bool {
			return c.FullPath() != route.METRICS
		}),
	))

	r.GET(route.METRICS, gin.WrapH(metricsHandler))

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	userGroup := r.Group("")
	{
		userGroup.POST(route.SIGN_UP, userHandler.SignUp)
		userGroup.POST(route.SIGN_IN, userHandler.SignIn)

		privateGroup := userGroup.Group("")
		{
			privateGroup.GET(route.ME, userHandler.GetMe)
		}

		nonAPIGatewayGroup := userGroup.Group("")
		{
			nonAPIGatewayGroup.GET(route.EMAILS, userHandler.GetEmailsByUserIDs) // REMOVE after switch to gRPC
		}
	}

	return r
}
