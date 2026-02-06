package middleware

import (
	httpInfrastructure "github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/transport/http"
	httpContextUser "github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/transport/http/context/user"

	"github.com/gin-gonic/gin"
)

func AdminRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, ok := httpContextUser.Get(ctx)
		if !ok {
			httpInfrastructure.NewUnauthorizedError(ctx)
			ctx.Abort()
			return
		}

		if !user.IsAdmin {
			httpInfrastructure.NewForbiddenError(ctx)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
