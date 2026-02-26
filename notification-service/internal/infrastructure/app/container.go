package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	sharedKafka "github.com/Yarik7610/library-backend-common/broker/kafka"
	"golang.org/x/sync/errgroup"

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
	httpServer           *http.Server
	stopOnce             sync.Once
	shutdownTracing      func(context.Context) error
}

func NewContainer() *Container {
	config, err := config.Parse()
	if err != nil {
		log.Fatalf("Config parse error: %v\n", err)
	}

	logger := logging.NewLogger(config.Env)

	shutdownTracing, err := tracing.Init(config)
	if err != nil {
		logger.Fatal(context.Background(), "Tracing init error", logging.Error(err))
	}

	subscriptionClient := subscription.NewClient()

	bookAddedReader := kafka.NewOtelReader(config, sharedKafka.BOOK_ADDED_TOPIC, sharedKafka.BOOK_ADDED_CONSUMER_GROUP_ID)
	bookAddedEmailSender := email.NewSender(config.Mail, config.MailPassword)
	bookAddedEmailSender.WithSubject("Book category subscription notification")

	bookAddedNotificator := bookadded.NewNotificator(logger, bookAddedReader, bookAddedEmailSender, subscriptionClient)

	httpServer, err := newHTTPServer(config)
	if err != nil {
		logger.Fatal(context.Background(), "HTTP server init error", logging.Error(err))
	}

	return &Container{
		Config:               config,
		Logger:               logger,
		bookAddedNotificator: bookAddedNotificator,
		httpServer:           httpServer,
		shutdownTracing:      shutdownTracing,
	}
}

func (c *Container) Start() error {
	group, ctx := errgroup.WithContext(context.Background())

	group.Go(func() error {
		c.bookAddedNotificator.Run(ctx)
		return nil
	})

	group.Go(func() error {
		err := c.httpServer.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	})

	// If HTTP server returns err, group.Wait() will still wait until all Goroutines are done.
	// But c.bookAddedNotificator.Run() is an infinite loop, so we won't catch error.
	// Thus, i create errgroup with context, and handle its cancellation in another Goroutine.

	// Context is canceled the first time a function passed to Goroutine returns a non-nil error
	// or the first time Wait returns,
	// whichever occurs first.
	group.Go(func() error {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		return c.Stop(shutdownCtx)
	})

	return group.Wait()
}

func (c *Container) Stop(ctx context.Context) error {
	var stopErr error

	c.stopOnce.Do(func() {
		c.bookAddedNotificator.Stop(ctx)

		if err := c.httpServer.Shutdown(ctx); err != nil {
			stopErr = err
			return
		}

		if err := c.shutdownTracing(ctx); err != nil {
			stopErr = err
			return
		}
	})

	return stopErr
}
