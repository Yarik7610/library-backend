package router

import (
	"github.com/Yarik7610/library-backend-common/transport/http/route"
	"github.com/Yarik7610/library-backend/api-gateway/internal/app/middleware"
	"github.com/Yarik7610/library-backend/api-gateway/internal/core"
	"github.com/gin-gonic/gin"
)

func registerUserRoutes(r *gin.Engine, userMicroserviceHandler gin.HandlerFunc) {
	userGroup := r.Group("")
	{
		userGroup.POST(route.SIGN_UP, userMicroserviceHandler)
		userGroup.POST(route.SIGN_IN, userMicroserviceHandler)

		privateGroup := userGroup.Group("")
		privateGroup.Use(middleware.AuthRequired(), core.InjectHeaders())
		{
			privateGroup.GET(route.ME, userMicroserviceHandler)
		}
	}
}
