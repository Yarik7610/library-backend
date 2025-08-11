package controller

import (
	"github.com/Yarik7610/library-backend/catalog-service/internal/service"
	"github.com/gin-gonic/gin"
)

type BookController interface {
	PreviewBook(ctx *gin.Context)
}

type bookController struct {
	bookService service.BookService
}

func NewBookController(bookService service.BookService) BookController {
	return &bookController{bookService: bookService}
}

func (c *bookController) PreviewBook(ctx *gin.Context) {

}
