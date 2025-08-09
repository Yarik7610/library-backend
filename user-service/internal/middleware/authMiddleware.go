package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Yarik7610/library-backend/user-service/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !utils.IsPrivateRoute(ctx.Request.URL.Path) {
			ctx.Next()
			return
		}

		authHeader := ctx.Request.Header.Get("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format, use Bearer <token>"})
			ctx.Abort()
			return
		}

		claims, err := utils.VerifyJWTToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		userID, _ := strconv.ParseUint(claims.Subject, 10, 64)
		isAdmin := false
		if len(claims.Audience) > 0 {
			isAdmin, _ = strconv.ParseBool(claims.Audience[0])
		}

		ctx.Set("userID", uint(userID))
		ctx.Set("isAdmin", isAdmin)

		zap.S().Infow("Authenticated", "userID", userID, "isAdmin", isAdmin)
		ctx.Next()
	}
}
