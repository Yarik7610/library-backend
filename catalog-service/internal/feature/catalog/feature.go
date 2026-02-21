package catalog

import (
	postgresRepositories "github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/postgres"
	redisRepositories "github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/redis"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/service"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/transport/http"
	kafkaInfrastructure "github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/broker/kafka"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/observability/logging"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/observability/metrics"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/storage/postgres/seed"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Feature struct {
	HTTPRouter *gin.Engine
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
	httpCatalogHandler := http.NewCatalogHandler(config, logger, catalogService)
	metricsHandler, err := metrics.Init()
	if err != nil {
		return nil, err
	}

	httpRouter := http.NewRouter(config, metricsHandler, httpCatalogHandler)
	return &Feature{HTTPRouter: httpRouter}, nil
}
