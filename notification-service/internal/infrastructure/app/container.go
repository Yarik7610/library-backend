package app

import (
	"context"
	"log"

	sharedKafka "github.com/Yarik7610/library-backend-common/broker/kafka"

	"github.com/Yarik7610/library-backend/notification-service/internal/core/notificator/bookadded"
	"github.com/Yarik7610/library-backend/notification-service/internal/infrastructure/broker/kafka"
	"github.com/Yarik7610/library-backend/notification-service/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/notification-service/internal/infrastructure/email"
	"github.com/Yarik7610/library-backend/notification-service/internal/infrastructure/observability/logging"
	"github.com/Yarik7610/library-backend/notification-service/internal/infrastructure/observability/tracing"
	"github.com/Yarik7610/library-backend/notification-service/internal/infrastructure/transport/http/microservice/subscription"
)

type Container struct {
	Config               *config.Config
	Logger               *logging.Logger
	bookAddedNotificator bookadded.Notificator
	shutdownTracing      func(context.Context) error
}

func NewContainer() *Container {
	config, err := config.Parse()
	if err != nil {
		log.Fatalf("Config load error: %v\n", err)
	}

	logger := logging.NewLogger(config.Env)

	shutdownTracing, err := tracing.Init(config)
	if err != nil {
		logger.Fatal(context.Background(), "Tracing init error", logging.Error(err))
	}

	subscriptionClient := subscription.NewClient()

	bookAddedReader := kafka.NewReader(sharedKafka.BOOK_ADDED_TOPIC, sharedKafka.BOOK_ADDED_CONSUMER_GROUP_ID)
	bookAddedEmailSender := email.NewSender(config.Mail, config.MailPassword)
	bookAddedEmailSender.WithSubject("Book category subscription notification")

	bookAddedNotificator := bookadded.NewNotificator(logger, bookAddedReader, bookAddedEmailSender, subscriptionClient)

	return &Container{
		Config:               config,
		Logger:               logger,
		bookAddedNotificator: bookAddedNotificator,
		shutdownTracing:      shutdownTracing,
	}
}

func (c *Container) Start() error {
	c.bookAddedNotificator.Run()
	return nil
}

func (c *Container) Stop(ctx context.Context) error {
	c.bookAddedNotificator.Stop()
	return nil
}
