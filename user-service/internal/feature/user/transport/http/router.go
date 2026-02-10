package http

import (
	"github.com/Yarik7610/library-backend-common/transport/http/route"
	"github.com/Yarik7610/library-backend/user-service/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(userHandler UserHandler) *gin.Engine {
	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST(route.SIGN_UP, userHandler.SignUp)
	r.POST(route.SIGN_IN, userHandler.SignIn)
	r.GET(route.ME, userHandler.GetMe)

	nonAPIGatewayGroup := r.Group("")
	{
		nonAPIGatewayGroup.GET(route.EMAILS, userHandler.GetEmailsByUserIDs)
	}

	return r
}
