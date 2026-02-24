package middleware

import (
	userContext "github.com/Yarik7610/library-backend/api-gateway/internal/app/context/user"
	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/errs"
	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/observability/tracing"

	httpInfrastructure "github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/transport/http"
	"github.com/gin-gonic/gin"

	"go.opentelemetry.io/otel/trace"
)

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		span := trace.SpanFromContext(c.Request.Context())

		user, ok := userContext.Get(c)
		if !ok {
			err := errs.NewUnauthorizedError()
			tracing.Error(span, err)
			httpInfrastructure.RenderError(c, err)
			c.Abort()
			return
		}

		if !user.IsAdmin {
			err := errs.NewForbiddenError()
			tracing.Error(span, err)
			httpInfrastructure.RenderError(c, err)
			c.Abort()
			return
		}

		c.Next()
	}
}
