package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewUnauthorizedError(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, "The token is missing, invalid or expired")

}

func NewForbiddenError(ctx *gin.Context) {
	ctx.JSON(http.StatusForbidden, "The token is valid, but lacks permission")
}
