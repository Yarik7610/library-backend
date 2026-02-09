package context

import (
	"strconv"

	"github.com/Yarik7610/library-backend-common/sharedconstants"
	"github.com/gin-gonic/gin"
)

func GetUserID(ctx *gin.Context) uint64 {
	userIDString := ctx.GetHeader(sharedconstants.HEADER_USER_ID)
	userID, err := strconv.ParseUint(userIDString, 10, 64)
	if err != nil {
		userID = 0
	}
	return userID
}
