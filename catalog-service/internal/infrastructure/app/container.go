package app

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	sharedKafka "github.com/Yarik7610/library-backend-common/broker/kafka"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/broker/kafka"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/observability/logging"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/observability/tracing"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/storage/postgres"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/storage/redis"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type Container struct {
	Config          *config.Config
	Logger          *logging.Logger
	httpServer      *http.Server
	gRPCServer      *grpc.Server
	stopOnce        sync.Once
	shutdownTracing func(context.Context) error
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

	postgresDB, err := postgres.Connect(config)
	if err != nil {
		logger.Fatal(context.Background(), "Postgres connect error", logging.Error(err))
	}

	redisClient, err := redis.Connect(config)
	if err != nil {
		logger.Fatal(context.Background(), "Redis connect error", logging.Error(err))
	}

	bookAddedWriter := kafka.NewOtelWriter(config, sharedKafka.BOOK_ADDED_TOPIC)

	catalogFeature, err := catalog.NewFeature(config, logger, postgresDB, redisClient, bookAddedWriter)
	if err != nil {
		logger.Fatal(context.Background(), "Catalog feature init error", logging.Error(err))
	}

	return &Container{
		Config:          config,
		Logger:          logger,
		httpServer:      catalogFeature.HTTPServer,
		gRPCServer:      catalogFeature.GRPCServer,
		shutdownTracing: shutdownTracing,
	}
}

func (c *Container) Start() error {
	group, ctx := errgroup.WithContext(context.Background())

	group.Go(func() error {
		err := c.httpServer.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	})

	group.Go(func() error {
		listener, err := net.Listen("tcp", ":"+c.Config.GRPCServerPort)
		if err != nil {
			return err
		}
		return c.gRPCServer.Serve(listener)
	})

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
		if err := c.httpServer.Shutdown(ctx); err != nil {
			stopErr = err
			return
		}

		c.gRPCServer.GracefulStop()

		if err := c.shutdownTracing(ctx); err != nil {
			stopErr = err
			return
		}
	})

	return stopErr
}
