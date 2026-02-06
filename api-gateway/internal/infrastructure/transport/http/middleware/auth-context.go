package middleware

import (
	"strconv"
	"strings"

	"github.com/Yarik7610/library-backend/api-gateway/config"
	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/jwt"
	httpInfrastructure "github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/transport/http"

	httpContextUser "github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/transport/http/context/user"
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
			httpInfrastructure.NewUnauthorizedError(ctx)
			ctx.Abort()
			return
		}

		claims, err := jwt.Verify(tokenString, config.Data.JWTSecret)
		if err != nil {
			httpInfrastructure.NewUnauthorizedError(ctx)
			ctx.Abort()
			return
		}

		userID, _ := strconv.ParseUint(claims.Subject, 10, 64)
		isAdmin := false
		if len(claims.Audience) > 0 {
			isAdmin, _ = strconv.ParseBool(claims.Audience[0])
		}

		httpContextUser.Set(ctx, httpContextUser.User{
			ID:      userID,
			IsAdmin: isAdmin,
		})

		ctx.Next()
	}
}
