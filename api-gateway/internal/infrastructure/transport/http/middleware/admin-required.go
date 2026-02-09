package middleware

import (
	httpInfrastructure "github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/transport/http"
	httpUserContext "github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/transport/http/context/user"

	"github.com/gin-gonic/gin"
)

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := httpUserContext.Get(c)
		if !ok {
			httpInfrastructure.NewUnauthorizedError(c)
			c.Abort()
			return
		}

		if !user.IsAdmin {
			httpInfrastructure.NewForbiddenError(c)
			c.Abort()
			return
		}

		c.Next()
	}
}
