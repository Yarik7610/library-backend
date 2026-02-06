package middleware

import (
	"strconv"

	"github.com/Yarik7610/library-backend-common/sharedconstants"
	"github.com/gin-gonic/gin"

	httpContextUser "github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/transport/http/context/user"
)

func InjectHeaders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, ok := httpContextUser.Get(ctx)
		if ok {
			ctx.Request.Header.Set(sharedconstants.HEADER_USER_ID, strconv.FormatUint(user.ID, 10))
			ctx.Request.Header.Set(sharedconstants.HEADER_IS_ADMIN, strconv.FormatBool(user.IsAdmin))
		}
		ctx.Next()
	}
}
