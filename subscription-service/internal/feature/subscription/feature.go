package subscription

import (
	"net/http"

	pb "github.com/Yarik7610/library-backend-common/transport/grpc/microservice/subscription"
	"github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription/repository/postgres"
	"github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription/service"
	grpcTransport "github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription/transport/grpc"
	httpTransport "github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription/transport/http"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/observability/logging"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/observability/metrics"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/transport/grpc/client/catalog"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/transport/grpc/client/user"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type Feature struct {
	HTTPServer *http.Server
	GRPCServer *grpc.Server
}

func NewFeature(
	config *config.Config,
	logger *logging.Logger,
	postgresDB *gorm.DB,
	catalogMicroserviceClient catalog.Client,
	userMicroserviceClient user.Client,
) (*Feature, error) {
	userBookCategorySubscriptionRepository := postgres.NewUserBookCategorySubscriptionRepository(postgresDB)

	subscriptionService := service.NewSubscriptionService(userBookCategorySubscriptionRepository, catalogMicroserviceClient, userMicroserviceClient)

	metricsHandler, err := metrics.Init()
	if err != nil {
		return nil, err
	}
	httpSubscriptionHandler := httpTransport.NewSubscriptionHandler(config, logger, subscriptionService)
	grpcSubscriptionHandler := grpcTransport.NewSubscriptionHandler(config, logger, subscriptionService)

	httpRouter := httpTransport.NewRouter(config, metricsHandler, httpSubscriptionHandler)
	httpServer := &http.Server{
		Addr:    ":" + config.HTTPServerPort,
		Handler: httpRouter,
	}

	gRPCServer := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
	pb.RegisterSubscriptionServiceServer(gRPCServer, grpcSubscriptionHandler)

	return &Feature{HTTPServer: httpServer, GRPCServer: gRPCServer}, nil
}
