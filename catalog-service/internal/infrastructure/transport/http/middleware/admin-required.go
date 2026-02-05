package middleware

import (
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/errs"
	httpInfrastructure "github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/transport/http"

	httpContext "github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/transport/http/context"

	"github.com/gin-gonic/gin"
)

func AdminRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		isAdminAny, ok := ctx.Get(httpContext.KeyIsAdmin)
		if !ok {
			httpInfrastructure.RenderError(ctx, errs.NewUnauthorizedError())
			ctx.Abort()
			return
		}

		isAdmin, ok := isAdminAny.(bool)
		if !ok || !isAdmin {
			httpInfrastructure.RenderError(ctx, errs.NewForbiddenError())
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
