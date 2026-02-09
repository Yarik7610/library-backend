package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Error string `json:"error"`
}

func NewUnauthorizedError(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Error{Error: "The token is missing, invalid or expired"})

}

func NewForbiddenError(c *gin.Context) {
	c.JSON(http.StatusForbidden, Error{Error: "The token is valid, but lacks permission"})
}
