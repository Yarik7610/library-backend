package middleware

import (
	"strconv"
	"strings"

	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/jwt"
	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/observability/tracing"
	"go.opentelemetry.io/otel/trace"

	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/errs"

	userContext "github.com/Yarik7610/library-backend/api-gateway/internal/app/context/user"
	httpInfrastructure "github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/transport/http"
	"github.com/gin-gonic/gin"
)

func AuthContext(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		span := trace.SpanFromContext(c.Request.Context())

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			err := errs.NewUnauthorizedError()
			tracing.Error(span, err)
			httpInfrastructure.RenderError(c, err)
			c.Abort()
			return
		}

		claims, err := jwt.Verify(tokenString, config.JWTSecret)
		if err != nil {
			err := errs.NewUnauthorizedError().WithCause(err)
			tracing.Error(span, err)
			httpInfrastructure.RenderError(c, err)
			c.Abort()
			return
		}

		userID, _ := strconv.ParseUint(claims.Subject, 10, 64)
		isAdmin := false
		if len(claims.Audience) > 0 {
			isAdmin, _ = strconv.ParseBool(claims.Audience[0])
		}

		userContext.Set(c, userContext.User{
			ID:      userID,
			IsAdmin: isAdmin,
		})

		c.Next()
	}
}
