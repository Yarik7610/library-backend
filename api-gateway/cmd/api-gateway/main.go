package main

import (
	"github.com/Yarik7610/library-backend/api-gateway/internal/app"
	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/observability/logging"
)

func main() {
	container := app.NewContainer()

	if err := container.Router.Run(":" + container.Config.ServerPort); err != nil {
		container.Logger.Fatal("Start server error", logging.String("port", container.Config.ServerPort), logging.Error(err))
	}
}
