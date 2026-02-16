package middleware

import (
	userContext "github.com/Yarik7610/library-backend/api-gateway/internal/app/context/user"
	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/errs"

	"github.com/gin-gonic/gin"
)

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := userContext.Get(c)
		if !ok {
			errs.NewUnauthorizedError(c)
			c.Abort()
			return
		}

		if !user.IsAdmin {
			errs.NewForbiddenError(c)
			c.Abort()
			return
		}

		c.Next()
	}
}
