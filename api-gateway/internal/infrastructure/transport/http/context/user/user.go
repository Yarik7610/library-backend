package context

import "github.com/gin-gonic/gin"

type User struct {
	ID      uint64
	IsAdmin bool
}

const userKey = "user"

func Set(ctx *gin.Context, user User) {
	ctx.Set(userKey, user)
}

func Get(ctx *gin.Context) (User, bool) {
	val, ok := ctx.Get(userKey)
	if !ok {
		return User{}, false
	}
	user, ok := val.(User)
	return user, ok
}
