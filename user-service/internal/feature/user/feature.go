package user

import (
	"github.com/Yarik7610/library-backend/user-service/internal/feature/user/repository/postgres"
	"github.com/Yarik7610/library-backend/user-service/internal/feature/user/service"
	"github.com/Yarik7610/library-backend/user-service/internal/feature/user/transport/http"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/storage/postgres/seed"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Feature struct {
	HTTPRouter *gin.Engine
}

func NewFeature(postgresDB *gorm.DB) *Feature {
	userRepository := postgres.NewUserRepository(postgresDB)

	seed.Admin(userRepository)

	userService := service.NewUserService(userRepository)
	userHandler := http.NewUserHandler(userService)

	httpRouter := http.NewRouter(userHandler)
	return &Feature{HTTPRouter: httpRouter}
}
