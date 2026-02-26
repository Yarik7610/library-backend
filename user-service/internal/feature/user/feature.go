package user

import (
	"net/http"

	pb "github.com/Yarik7610/library-backend-common/transport/grpc/microservice/user"
	"github.com/Yarik7610/library-backend/user-service/internal/feature/user/repository/postgres"
	"github.com/Yarik7610/library-backend/user-service/internal/feature/user/service"
	grpcTransport "github.com/Yarik7610/library-backend/user-service/internal/feature/user/transport/grpc"
	httpTransport "github.com/Yarik7610/library-backend/user-service/internal/feature/user/transport/http"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/observability/logging"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/observability/metrics"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/storage/postgres/seed"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type Feature struct {
	HTTPServer *http.Server
	GRPCServer *grpc.Server
}

func NewFeature(config *config.Config, logger *logging.Logger, postgresDB *gorm.DB) (*Feature, error) {
	userRepository := postgres.NewUserRepository(postgresDB)

	if err := seed.Admin(config, userRepository); err != nil {
		return nil, err
	}

	userService := service.NewUserService(config, userRepository)

	metricsHandler, err := metrics.Init()
	if err != nil {
		return nil, err
	}
	httpUserHandler := httpTransport.NewUserHandler(config, logger, userService)
	grpcUserHandler := grpcTransport.NewUserHandler(config, logger, userService)

	httpRouter := httpTransport.NewRouter(config, metricsHandler, httpUserHandler)
	httpServer := &http.Server{
		Addr:    ":" + config.HTTPServerPort,
		Handler: httpRouter,
	}

	gRPCServer := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
	pb.RegisterUserServiceServer(gRPCServer, grpcUserHandler)

	return &Feature{HTTPServer: httpServer, GRPCServer: gRPCServer}, nil
}
