package middleware

import (
	userContext "github.com/Yarik7610/library-backend/api-gateway/internal/app/context/user"
	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/errs"
	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/observability/tracing"
	httpInfrastructure "github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/transport/http"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := userContext.Get(c); !ok {
			span := trace.SpanFromContext(c.Request.Context())

			err := errs.NewUnauthorizedError()
			tracing.Error(span, err)
			httpInfrastructure.RenderError(c, err)
			c.Abort()
			return
		}
		c.Next()
	}
}
