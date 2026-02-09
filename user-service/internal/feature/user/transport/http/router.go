package http

import (
	"github.com/Yarik7610/library-backend-common/sharedconstants"
	"github.com/Yarik7610/library-backend/user-service/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(userHandler UserHandler) *gin.Engine {
	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST(sharedconstants.SIGN_UP_ROUTE, userHandler.SignUp)
	r.POST(sharedconstants.SIGN_IN_ROUTE, userHandler.SignIn)
	r.GET(sharedconstants.ME_ROUTE, userHandler.GetMe)

	nonAPIGatewayGroup := r.Group("")
	{
		nonAPIGatewayGroup.GET(sharedconstants.EMAILS_ROUTE, userHandler.GetEmailsByUserIDs)
	}

	return r
}
