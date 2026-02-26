package catalog

import (
	"net/http"

	pb "github.com/Yarik7610/library-backend-common/transport/grpc/microservice/catalog"
	postgresRepositories "github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/postgres"
	redisRepositories "github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/redis"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/service"
	grpcTransport "github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/transport/grpc"
	httpTransport "github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/transport/http"
	kafkaInfrastructure "github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/broker/kafka"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/observability/logging"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/observability/metrics"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/storage/postgres/seed"
	"github.com/redis/go-redis/v9"
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
	redisClient *redis.Client,
	bookAddedWriter *kafkaInfrastructure.OtelWriter,
) (*Feature, error) {
	redisBookRepository := redisRepositories.NewBookRepository(redisClient)
	postgresBookRepository := postgresRepositories.NewBookRepository(postgresDB)
	postgresPageRepository := postgresRepositories.NewPageRepository(postgresDB)
	postgresAuthorRepository := postgresRepositories.NewAuthorRepository(postgresDB)

	if err := seed.Books(postgresBookRepository, postgresPageRepository, postgresAuthorRepository); err != nil {
		return nil, err
	}

	catalogService := service.NewCatalogService(
		logger, postgresDB,
		bookAddedWriter, redisBookRepository,
		postgresAuthorRepository, postgresBookRepository, postgresPageRepository,
	)

	metricsHandler, err := metrics.Init()
	if err != nil {
		return nil, err
	}
	httpCatalogHandler := httpTransport.NewCatalogHandler(config, logger, catalogService)
	gRPCCatalogHandler := grpcTransport.NewCatalogHandler(config, logger, catalogService)

	httpRouter := httpTransport.NewRouter(config, metricsHandler, httpCatalogHandler)
	httpServer := &http.Server{
		Addr:    ":" + config.HTTPServerPort,
		Handler: httpRouter,
	}

	gRPCServer := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
	pb.RegisterCatalogServiceServer(gRPCServer, gRPCCatalogHandler)

	return &Feature{HTTPServer: httpServer, GRPCServer: gRPCServer}, nil
}
