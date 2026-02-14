package main

import (
	"github.com/Yarik7610/library-backend/notification-service/internal/infrastructure/app"
)

func main() {
	container := app.NewContainer()

	container.BookAddedNotificator.Run()
	defer container.BookAddedNotificator.Stop()
}
