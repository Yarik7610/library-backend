package middleware

import (
	"strconv"

	"github.com/Yarik7610/library-backend-common/transport/http/header"
	"github.com/gin-gonic/gin"

	httpUserContext "github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/transport/http/context/user"
)

func InjectHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := httpUserContext.Get(c)
		if ok {
			c.Request.Header.Set(header.USER_ID, strconv.FormatUint(user.ID, 10))
			c.Request.Header.Set(header.IS_ADMIN, strconv.FormatBool(user.IsAdmin))
		}
		c.Next()
	}
}
