package middleware

import (
	httpInfrastructure "github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/transport/http"
	httpContextUser "github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/transport/http/context/user"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if _, ok := httpContextUser.Get(ctx); !ok {
			httpInfrastructure.NewUnauthorizedError(ctx)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
