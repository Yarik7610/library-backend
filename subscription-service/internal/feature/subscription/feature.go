package subscription

import (
	"github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription/repository/postgres"
	"github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription/service"
	"github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription/transport/http"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/observability/logging"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/observability/metrics"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/transport/http/microservice/catalog"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/transport/http/microservice/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Feature struct {
	HTTPRouter *gin.Engine
}

func NewFeature(config *config.Config, logger *logging.Logger, postgresDB *gorm.DB) (*Feature, error) {
	catalogMicroserviceClient := catalog.NewClient()
	userMicroserviceClient := user.NewClient()

	userBookCategorySubscriptionRepository := postgres.NewUserBookCategorySubscriptionRepository(postgresDB)
	subscriptionService := service.NewSubscriptionService(userBookCategorySubscriptionRepository, catalogMicroserviceClient, userMicroserviceClient)
	httpSubscriptionHandler := http.NewSubscriptionHandler(logger, subscriptionService)
	metricsHandler, err := metrics.Init()
	if err != nil {
		return nil, err
	}

	httpRouter := http.NewRouter(config, metricsHandler, httpSubscriptionHandler)
	return &Feature{HTTPRouter: httpRouter}, nil
}
