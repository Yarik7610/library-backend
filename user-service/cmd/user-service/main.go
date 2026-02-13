package main

import (
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/app"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/observability/logging"
)

func main() {
	container := app.NewContainer()

	if err := container.UserFeature.HTTPRouter.Run(":" + container.Config.ServerPort); err != nil {
		container.Logger.Fatal("Start server error", logging.String("port", container.Config.ServerPort), logging.Error(err))
	}
}
