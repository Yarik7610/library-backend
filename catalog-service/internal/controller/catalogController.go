package controller

import (
	"net/http"
	"strconv"

	"github.com/Yarik7610/library-backend/catalog-service/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CatalogController interface {
	PreviewBook(ctx *gin.Context)
	GetCategories(ctx *gin.Context)
	GetAuthorsBooks(ctx *gin.Context)
}

type catalogController struct {
	catalogService service.CatalogService
}

func NewCatalogController(catalogService service.CatalogService) CatalogController {
	return &catalogController{catalogService: catalogService}
}

func (c *catalogController) PreviewBook(ctx *gin.Context) {
	bookIDStr := ctx.Param("bookID")
	bookID, err := strconv.ParseUint(bookIDStr, 10, 64)

	if err != nil {
		zap.S().Errorf("Preview book id param error: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, customErr := c.catalogService.PreviewBook(uint(bookID))
	if customErr != nil {
		zap.S().Errorf("Preview book error: %v\n", err)
		ctx.JSON(customErr.Code, gin.H{"error": customErr.Message})
		return
	}

	ctx.JSON(http.StatusOK, book)
}

func (c *catalogController) GetCategories(ctx *gin.Context) {
	categories, err := c.catalogService.GetCategories()
	if err != nil {
		zap.S().Errorf("List categories error: %v\n", err)
		ctx.JSON(err.Code, gin.H{"error": err.Message})
		return
	}

	ctx.JSON(http.StatusOK, categories)
}

func (c *catalogController) GetAuthorsBooks(ctx *gin.Context) {
	authorName := ctx.Param("authorName")
	authorsBooks, err := c.catalogService.GetAuthorsBooks(authorName)
	if err != nil {
		zap.S().Errorf("Get authors books error: %v\n", err)
		ctx.JSON(err.Code, gin.H{"error": err.Message})
		return
	}

	ctx.JSON(http.StatusOK, authorsBooks)
}
