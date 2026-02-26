package swagger

import (
	"net/http"

	"github.com/Yarik7610/library-backend-common/microservice"
	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/observability/logging"
	httpInfrastructure "github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/transport/http"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	GetMergedDocs(c *gin.Context)
}

type handler struct {
	config *config.Config
	logger *logging.Logger
}

func NewHandler(config *config.Config, logger *logging.Logger) Handler {
	return &handler{config: config, logger: logger}
}

func (h *handler) GetMergedDocs(c *gin.Context) {
	ctx := c.Request.Context()

	userDocs, err := fetchDocsJSON(microservice.USER_HTTP_ADDRESS)
	if err != nil {
		h.logger.Error(ctx, "User microservice swagger fetch error", logging.Error(err))
		httpInfrastructure.RenderError(c, err)
		return
	}

	catalogDocs, err := fetchDocsJSON(microservice.CATALOG_HTTP_ADDRESS)
	if err != nil {
		h.logger.Error(ctx, "Catalog microservice swagger fetch error", logging.Error(err))
		httpInfrastructure.RenderError(c, err)
		return
	}

	subscriptionDocs, err := fetchDocsJSON(microservice.SUBSCRIPTIONS_HTTP_ADDRESS)
	if err != nil {
		h.logger.Error(ctx, "Subscription microservice swagger fetch error", logging.Error(err))
		httpInfrastructure.RenderError(c, err)
		return
	}

	mergedDoc := mergeDocs(userDocs, catalogDocs, subscriptionDocs)
	c.JSON(http.StatusOK, mergedDoc)
}
