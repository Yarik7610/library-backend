package subscription

import (
	"github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription/repository/postgres"
	"github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription/service"
	"github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription/transport/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Feature struct {
	HTTPRouter *gin.Engine
}

func NewFeature(postgresDB *gorm.DB) *Feature {
	userBookCategoryRepository := postgres.NewUserBookCategoryRepository(postgresDB)
	subscriptionService := service.NewSubscriptionService(userBookCategoryRepository)
	httpSubscriptionHandler := http.NewSubscriptionHandler(subscriptionService)

	httpRouter := http.NewRouter(httpSubscriptionHandler)
	return &Feature{HTTPRouter: httpRouter}
}
