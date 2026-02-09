package middleware

import (
	httpInfrastructure "github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/transport/http"
	httpUserContext "github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/transport/http/context/user"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := httpUserContext.Get(c); !ok {
			httpInfrastructure.NewUnauthorizedError(c)
			c.Abort()
			return
		}
		c.Next()
	}
}
