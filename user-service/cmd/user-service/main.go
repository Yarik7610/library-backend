package main

import (
	"github.com/Yarik7610/library-backend/user-service/internal/feature/user"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/storage/postgres"
	"go.uber.org/zap"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
}

func main() {
	if err := config.Init(); err != nil {
		zap.S().Fatalf("Config load error: %v\n", err)
	}

	postgresDB := postgres.DB()

	userFeature := user.NewFeature(postgresDB)

	if err := userFeature.HTTPRouter.Run(":" + config.Data.ServerPort); err != nil {
		zap.S().Fatalf("Start error on port %s: %v", config.Data.ServerPort, err)
	}
}
