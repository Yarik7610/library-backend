package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Error string `json:"error"`
}

func NewUnauthorizedError(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, Error{Error: "The token is missing, invalid or expired"})

}

func NewForbiddenError(ctx *gin.Context) {
	ctx.JSON(http.StatusForbidden, Error{Error: "The token is valid, but lacks permission"})
}
