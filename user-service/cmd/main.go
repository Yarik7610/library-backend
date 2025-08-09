package main

import (
	"github.com/Yarik7610/library-backend/user-service/config"
	"github.com/Yarik7610/library-backend/user-service/internal/constants"
	"github.com/Yarik7610/library-backend/user-service/internal/controller"
	"github.com/Yarik7610/library-backend/user-service/internal/middleware"
	"github.com/Yarik7610/library-backend/user-service/internal/model"
	"github.com/Yarik7610/library-backend/user-service/internal/repository"
	"github.com/Yarik7610/library-backend/user-service/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
}

func main() {
	err := config.Init()
	if err != nil {
		zap.S().Fatalf("Config load error: %v\n", err)
	}

	db, err := gorm.Open(postgres.Open(config.Data.PostgresURL), &gorm.Config{})
	if err != nil {
		zap.S().Fatalf("Gorm open error: %v\n", err)
	}
	zap.S().Info("Successfully connected to Postgres")

	err = db.AutoMigrate(&model.User{}, &model.BookCategory{})
	if err != nil {
		zap.S().Fatalf("Gorm auto migrate error: %v", err)
	}
	zap.S().Info("Successfully made auto migrate")

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	r := gin.Default()
	r.Use(middleware.AuthMiddleware())

	r.POST(constants.SIGN_UP_ROUTE, userController.SignUp)
	r.POST(constants.SIGN_IN_ROUTE, userController.SignIn)
	r.GET(constants.ME_ROUTE, userController.Me)

	if err := r.Run(":" + config.Data.ServerPort); err != nil {
		zap.S().Fatalf("Server start error on port %s: %v", config.Data.ServerPort, err)
	}
}
