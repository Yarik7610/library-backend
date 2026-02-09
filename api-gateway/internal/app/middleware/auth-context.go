package middleware

import (
	"strconv"
	"strings"

	"github.com/Yarik7610/library-backend/api-gateway/internal/app"
	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/jwt"

	userContext "github.com/Yarik7610/library-backend/api-gateway/internal/app/context/user"
	"github.com/gin-gonic/gin"
)

func AuthContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			app.NewUnauthorizedError(c)
			c.Abort()
			return
		}

		claims, err := jwt.Verify(tokenString, config.Data.JWTSecret)
		if err != nil {
			app.NewUnauthorizedError(c)
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
