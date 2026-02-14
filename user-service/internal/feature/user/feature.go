package user

import (
	"github.com/Yarik7610/library-backend/user-service/internal/feature/user/repository/postgres"
	"github.com/Yarik7610/library-backend/user-service/internal/feature/user/service"
	"github.com/Yarik7610/library-backend/user-service/internal/feature/user/transport/http"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/observability/logging"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/storage/postgres/seed"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Feature struct {
	HTTPRouter *gin.Engine
}

func NewFeature(config *config.Config, logger *logging.Logger, postgresDB *gorm.DB) *Feature {
	userRepository := postgres.NewUserRepository(postgresDB)

	if err := seed.Admin(config, userRepository); err != nil {
		logger.Fatal("Postgres admin seed error", logging.Error(err))
	}

	userService := service.NewUserService(config, userRepository)
	userHandler := http.NewUserHandler(config, logger, userService)

	httpRouter := http.NewRouter(userHandler)
	return &Feature{HTTPRouter: httpRouter}
}
