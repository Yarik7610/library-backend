package middleware

import (
	"github.com/Yarik7610/library-backend/api-gateway/internal/app"
	userContext "github.com/Yarik7610/library-backend/api-gateway/internal/app/context/user"

	"github.com/gin-gonic/gin"
)

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := userContext.Get(c)
		if !ok {
			app.NewUnauthorizedError(c)
			c.Abort()
			return
		}

		if !user.IsAdmin {
			app.NewForbiddenError(c)
			c.Abort()
			return
		}

		c.Next()
	}
}
