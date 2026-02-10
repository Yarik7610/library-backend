package middleware

import (
	"fmt"

	"github.com/Yarik7610/library-backend/api-gateway/internal/app"
	userContext "github.com/Yarik7610/library-backend/api-gateway/internal/app/context/user"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := userContext.Get(c); !ok {
			fmt.Println("TWEWEW")
			app.NewUnauthorizedError(c)
			c.Abort()
			return
		}
		c.Next()
	}
}
