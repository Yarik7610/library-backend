package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Yarik7610/library-backend-common/sharedconstants"
	"github.com/Yarik7610/library-backend/api-gateway/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fullPath := ctx.FullPath()
		method := ctx.Request.Method

		if !utils.IsPrivateRoute(method, fullPath) && !utils.IsAdminRoute(method, fullPath) {
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

		if utils.IsAdminRoute(method, fullPath) && !isAdmin {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "route allowed for admin only"})
			ctx.Abort()
			return
		}

		ctx.Request.Header.Set(sharedconstants.HEADER_USER_ID, fmt.Sprintf("%d", userID))
		ctx.Request.Header.Set(sharedconstants.HEADER_IS_ADMIN, strconv.FormatBool(isAdmin))

		zap.S().Infow("Authenticated", sharedconstants.HEADER_USER_ID, userID, sharedconstants.HEADER_IS_ADMIN, isAdmin)
		ctx.Next()
	}
}
