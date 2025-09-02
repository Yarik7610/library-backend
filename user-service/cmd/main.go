package main

import (
	"github.com/Yarik7610/library-backend-common/sharedconstants"
	"github.com/Yarik7610/library-backend/user-service/config"
	docs "github.com/Yarik7610/library-backend/user-service/docs"
	"github.com/Yarik7610/library-backend/user-service/internal/connect"
	"github.com/Yarik7610/library-backend/user-service/internal/controller"
	"github.com/Yarik7610/library-backend/user-service/internal/repository"
	"github.com/Yarik7610/library-backend/user-service/internal/seed"
	"github.com/Yarik7610/library-backend/user-service/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
}

func main() {
	err := config.Init()
	if err != nil {
		zap.S().Fatalf("Config load error: %v\n", err)
	}

	db := connect.DB()

	userRepo := repository.NewUserRepository(db)

	seed.Admin(userRepo)

	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST(sharedconstants.SIGN_UP_ROUTE, userController.SignUp)
	r.POST(sharedconstants.SIGN_IN_ROUTE, userController.SignIn)
	r.GET(sharedconstants.ME_ROUTE, userController.GetMe)

	nonAPIGatewayGroup := r.Group("")
	{
		nonAPIGatewayGroup.GET(sharedconstants.EMAILS_ROUTE, userController.GetEmailsByUserIDs)
	}

	if err := r.Run(":" + config.Data.ServerPort); err != nil {
		zap.S().Fatalf("Start error on port %s: %v", config.Data.ServerPort, err)
	}
}
