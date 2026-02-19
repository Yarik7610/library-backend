package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/app"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/observability/logging"
)

func main() {
	container := app.NewContainer()

	go func() {
		container.Logger.Info(context.Background(), "Starting container")
		if err := container.Start(); err != nil {
			container.Logger.Fatal(context.Background(),
				"Start container error",
				logging.String("HTTP port", container.Config.HTTPServerPort),
				logging.Error(err))
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	container.Logger.Info(context.Background(), "Shutting container down gracefully")
	if err := container.Stop(shutdownCtx); err != nil {
		container.Logger.Fatal(context.Background(), "Gracefull container shutdown failed")
	}
	container.Logger.Info(context.Background(), "Container stopped gracefully")
}
