package swagger

import (
	"net/http"

	"github.com/Yarik7610/library-backend-common/microservice"
	"github.com/Yarik7610/library-backend/api-gateway/docs"
	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/errs"
	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/observability/logging"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(logger *logging.Logger, r *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/"

	r.GET("/swagger-json/doc.json", func(c *gin.Context) {
		userDocs, err := fetchDocsJSON(microservice.USER_ADDRESS)
		if err != nil {
			logger.Error("User microservice swagger fetch error", logging.Error(err))
			errs.NewInternalServerError(c)
			return
		}

		catalogDocs, err := fetchDocsJSON(microservice.CATALOG_ADDRESS)
		if err != nil {
			logger.Error("Catalog microservice swagger fetch error", logging.Error(err))
			errs.NewInternalServerError(c)
			return
		}

		subscriptionDocs, err := fetchDocsJSON(microservice.SUBSCRIPTIONS_ADDRESS)
		if err != nil {
			logger.Error("Subscription microservice swagger fetch error", logging.Error(err))
			errs.NewInternalServerError(c)
			return
		}

		mergedDoc := mergeDocs(userDocs, catalogDocs, subscriptionDocs)
		c.JSON(http.StatusOK, mergedDoc)
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger-json/doc.json")))
}
