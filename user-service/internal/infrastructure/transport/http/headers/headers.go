package headers

import (
	"strconv"

	"github.com/Yarik7610/library-backend-common/transport/http/header"
	"github.com/gin-gonic/gin"
)

func GetUserID(ctx *gin.Context) (uint64, error) {
	userIDString := ctx.GetHeader(header.HEADER_USER_ID)
	userID, err := strconv.ParseUint(userIDString, 10, 64)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
