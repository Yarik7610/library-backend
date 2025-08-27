package main

import (
	"fmt"

	"github.com/Yarik7610/library-backend/api-gateway/config"
	"github.com/Yarik7610/library-backend/api-gateway/internal/email"
	"github.com/Yarik7610/library-backend/api-gateway/internal/utils"

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

	body := fmt.Sprintf("Hello! New books arrival in %q category", utils.Capitalize("horror"))

	sender := email.NewSender()
	sender.WithSubject("Subscription notification")
	sender.Send(body, []string{config.Data.Mail})
}
