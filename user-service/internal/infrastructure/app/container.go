package app

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/Yarik7610/library-backend/user-service/internal/feature/user"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/observability/logging"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/observability/tracing"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/storage/postgres"
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

	userFeature, err := user.NewFeature(config, logger, postgresDB)
	if err != nil {
		logger.Fatal(context.Background(), "User feature init error", logging.Error(err))
	}

	return &Container{
		Config:          config,
		Logger:          logger,
		httpServer:      userFeature.HTTPServer,
		gRPCServer:      userFeature.GRPCServer,
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
