package app

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/Yarik7610/library-backend-common/microservice"
	"github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/observability/logging"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/observability/tracing"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/storage/postgres"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/transport/grpc/client/catalog"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/transport/grpc/client/user"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type Container struct {
	Config                      *config.Config
	Logger                      *logging.Logger
	httpServer                  *http.Server
	gRPCServer                  *grpc.Server
	gRPCUserMicroserviceConn    *grpc.ClientConn
	gRPCCatalogMicroserviceConn *grpc.ClientConn
	stopOnce                    sync.Once
	shutdownTracing             func(context.Context) error
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

	userMicroserviceClient, gRPCUserMicroserviceConn, err := user.NewClient()
	if err != nil {
		logger.Fatal(context.Background(),
			"gRPC user microservice client connect error",
			logging.String("gRPC server address", microservice.USER_GRPC_ADDRESS),
			logging.Error(err),
		)
	}
	catalogMicroserviceClient, gRPCCatalogMicroserviceConn, err := catalog.NewClient()
	if err != nil {
		logger.Fatal(context.Background(),
			"gRPC catalog microservice client connect error",
			logging.String("gRPC server address", microservice.CATALOG_GRPC_ADDRESS),
			logging.Error(err),
		)
	}

	subscriptionFeature, err := subscription.NewFeature(
		config, logger, postgresDB,
		catalogMicroserviceClient, userMicroserviceClient,
	)
	if err != nil {
		logger.Fatal(context.Background(), "Subscription feature init error", logging.Error(err))
	}

	return &Container{
		Config:                      config,
		Logger:                      logger,
		httpServer:                  subscriptionFeature.HTTPServer,
		gRPCServer:                  subscriptionFeature.GRPCServer,
		gRPCUserMicroserviceConn:    gRPCUserMicroserviceConn,
		gRPCCatalogMicroserviceConn: gRPCCatalogMicroserviceConn,
		shutdownTracing:             shutdownTracing,
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
		c.gRPCServer.GracefulStop()

		if err := c.httpServer.Shutdown(ctx); err != nil {
			stopErr = err
			return
		}

		if err := c.gRPCUserMicroserviceConn.Close(); err != nil {
			stopErr = err
			return
		}

		if err := c.gRPCCatalogMicroserviceConn.Close(); err != nil {
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
