package context

import "github.com/gin-gonic/gin"

type contextKey string

const (
	KeyUserID  contextKey = "userID"
	KeyIsAdmin contextKey = "isAdmin"
)

func GetUserID(ctx *gin.Context) (uint64, bool) {
	val, exists := ctx.Get(KeyUserID)
	if !exists {
		return 0, false
	}

	userID, ok := val.(uint64)
	return userID, ok
}
