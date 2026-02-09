package main

import (
	"github.com/Yarik7610/library-backend-common/broker"
	"github.com/Yarik7610/library-backend-common/sharedconstants"
	"github.com/Yarik7610/library-backend/notification-service/config"
	"github.com/Yarik7610/library-backend/notification-service/internal/core"
	"github.com/Yarik7610/library-backend/notification-service/internal/email"

	"go.uber.org/zap"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
}

func main() {
	if err := config.Init(); err != nil {
		zap.S().Fatalf("Config load error: %v\n", err)
	}

	bookAddedReader := broker.NewReader(sharedconstants.BOOK_ADDED_TOPIC, sharedconstants.BOOK_ADDED_CONSUMER_GROUP_ID)
	sender := email.NewSender(config.Data.Mail, config.Data.MailPassword)
	sender.WithSubject("Subscription notification")

	notificator := core.NewNotificator(bookAddedReader, sender)
	notificator.Run()
}
