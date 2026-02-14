package catalog

import (
	postgresRepositories "github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/postgres"
	redisRepositories "github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/redis"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/service"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/transport/http"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/observability/logging"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/storage/postgres/seed"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

type Feature struct {
	HTTPRouter *gin.Engine
}

func NewFeature(
	logger *logging.Logger,
	postgresDB *gorm.DB,
	redisClient *redis.Client,
	bookAddedWriter *kafka.Writer,
) *Feature {
	redisBookRepository := redisRepositories.NewBookRepository(redisClient)
	postgresBookRepository := postgresRepositories.NewBookRepository(postgresDB)
	postgresPageRepository := postgresRepositories.NewPageRepository(postgresDB)
	postgresAuthorRepository := postgresRepositories.NewAuthorRepository(postgresDB)

	if err := seed.Books(postgresBookRepository, postgresPageRepository, postgresAuthorRepository); err != nil {
		logger.Fatal("Postgres books seed error", logging.Error(err))
	}

	catalogService := service.NewCatalogService(logger, postgresDB, bookAddedWriter, redisBookRepository, postgresAuthorRepository, postgresBookRepository, postgresPageRepository)
	httpCatalogHandler := http.NewCatalogHandler(logger, catalogService)

	httpRouter := http.NewRouter(httpCatalogHandler)
	return &Feature{HTTPRouter: httpRouter}
}
