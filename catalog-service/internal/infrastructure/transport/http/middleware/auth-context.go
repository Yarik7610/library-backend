package middleware

import (
	"strconv"
	"strings"

	"github.com/Yarik7610/library-backend-common/jwt"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/errs"
	httpInfrastructure "github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/transport/http"
	httpContext "github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/transport/http/context"
	"github.com/gin-gonic/gin"
)

func AuthContext() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.Request.Header.Get("Authorization")
		if authHeader == "" {
			ctx.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			httpInfrastructure.RenderError(ctx, errs.NewUnauthorizedError())
			ctx.Abort()
			return
		}

		claims, err := jwt.Verify(tokenString, config.Data.JWTSecret)
		if err != nil {
			httpInfrastructure.RenderError(ctx, errs.NewUnauthorizedError().WithCause(err))
			ctx.Abort()
			return
		}

		userID, _ := strconv.ParseUint(claims.Subject, 10, 64)
		isAdmin := false
		if len(claims.Audience) > 0 {
			isAdmin, _ = strconv.ParseBool(claims.Audience[0])
		}

		ctx.Set(httpContext.KeyUserID, userID)
		ctx.Set(httpContext.KeyIsAdmin, isAdmin)

		ctx.Next()
	}
}
