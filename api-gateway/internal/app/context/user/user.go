package user

import "github.com/gin-gonic/gin"

type User struct {
	ID      uint64
	IsAdmin bool
}

const userKey = "user"

func Set(c *gin.Context, user User) {
	c.Set(userKey, user)
}

func Get(c *gin.Context) (User, bool) {
	val, ok := c.Get(userKey)
	if !ok {
		return User{}, false
	}
	user, ok := val.(User)
	return user, ok
}
