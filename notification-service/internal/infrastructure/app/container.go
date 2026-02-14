package app

import (
	"log"

	sharedKafka "github.com/Yarik7610/library-backend-common/broker/kafka"

	"github.com/Yarik7610/library-backend/notification-service/internal/core/notificator/bookadded"
	"github.com/Yarik7610/library-backend/notification-service/internal/infrastructure/broker/kafka"
	"github.com/Yarik7610/library-backend/notification-service/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/notification-service/internal/infrastructure/email"
	"github.com/Yarik7610/library-backend/notification-service/internal/infrastructure/observability/logging"
	"github.com/Yarik7610/library-backend/notification-service/internal/infrastructure/transport/http/microservice/subscription"
)

type Container struct {
	Config               *config.Config
	Logger               *logging.Logger
	BookAddedNotificator bookadded.Notificator
}

func NewContainer() *Container {
	config, err := config.Init()
	if err != nil {
		log.Fatalf("Config load error: %v\n", err)
	}

	logger := logging.NewLogger(config.Env)

	bookAddedReader := kafka.NewReader(sharedKafka.BOOK_ADDED_TOPIC, sharedKafka.BOOK_ADDED_CONSUMER_GROUP_ID)
	bookAddedEmailSender := email.NewSender(config.Mail, config.MailPassword)
	bookAddedEmailSender.WithSubject("Book category subscription notification")
	subscriptionClient := subscription.NewClient()

	bookAddedNotificator := bookadded.NewNotificator(logger, bookAddedReader, bookAddedEmailSender, subscriptionClient)

	return &Container{
		Config:               config,
		Logger:               logger,
		BookAddedNotificator: bookAddedNotificator,
	}
}
