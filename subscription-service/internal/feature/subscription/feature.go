package subscription

import (
	"github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription/repository/postgres"
	"github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription/service"
	"github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription/transport/http"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/observability/logging"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/transport/http/microservice/catalog"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/transport/http/microservice/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Feature struct {
	HTTPRouter *gin.Engine
}

func NewFeature(logger *logging.Logger, postgresDB *gorm.DB) *Feature {
	catalogMicroserviceClient := catalog.NewClient()
	userMicroserviceClient := user.NewClient()

	userBookCategorySubscriptionRepository := postgres.NewUserBookCategorySubscriptionRepository(postgresDB)
	subscriptionService := service.NewSubscriptionService(userBookCategorySubscriptionRepository, catalogMicroserviceClient, userMicroserviceClient)
	httpSubscriptionHandler := http.NewSubscriptionHandler(logger, subscriptionService)

	httpRouter := http.NewRouter(httpSubscriptionHandler)
	return &Feature{HTTPRouter: httpRouter}
}
