package controller

import (
	"net/http"
	"strconv"

	"github.com/Yarik7610/library-backend-common/custom"
	"github.com/Yarik7610/library-backend/catalog-service/internal/dto"
	"github.com/Yarik7610/library-backend/catalog-service/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CatalogController interface {
	GetCategories(ctx *gin.Context)
	ListBooksByCategory(ctx *gin.Context)
	PreviewBook(ctx *gin.Context)
	GetBooksByAuthorID(ctx *gin.Context)
	SearchBooks(ctx *gin.Context)
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

func (c *catalogController) GetBooksByAuthorID(ctx *gin.Context) {
	authorIDString := ctx.Param("authorID")
	authorID, err := strconv.Atoi(authorIDString)
	if err != nil {
		zap.S().Errorf("Get books by author ID atoi error: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var customErr *custom.Err
	books, customErr := c.catalogService.GetBooksByAuthorID(authorID)
	if customErr != nil {
		zap.S().Errorf("Get books by author ID error: %v\n", err)
		ctx.JSON(customErr.Code, gin.H{"error": customErr.Message})
		return
	}

	ctx.JSON(http.StatusOK, books)
}

func (c *catalogController) SearchBooks(ctx *gin.Context) {
	authorName := ctx.Query("author")
	title := ctx.Query("title")

	if authorName == "" && title == "" {
		zap.S().Errorf("Search books error: can't have both query strings author and title empty")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Can't have both empty author and title"})
		return
	}

	var books []dto.Books
	var err *custom.Err

	if authorName != "" && title != "" {
		books, err = c.catalogService.GetBooksByAuthorNameAndTitle(authorName, title)
	} else if authorName != "" {
		books, err = c.catalogService.GetBooksByAuthorName(authorName)
	} else {
		books, err = c.catalogService.GetBooksByTitle(title)
	}

	if err != nil {
		zap.S().Errorf("Search books error: %v\n", err.Message)
		ctx.JSON(err.Code, gin.H{"error": err.Message})
		return
	}

	ctx.JSON(http.StatusOK, books)
}

type ListBooksByCategoryQuery struct {
	Page  int    `form:"page" binding:"required,min=1"`
	Count int    `form:"count" binding:"required,min=1,max=100"`
	Sort  string `form:"sort"`
	Order string `form:"order"`
}

func (c *catalogController) ListBooksByCategory(ctx *gin.Context) {
	categoryName := ctx.Param("categoryName")
	var query ListBooksByCategoryQuery

	if err := ctx.ShouldBindQuery(&query); err != nil {
		zap.S().Errorf("List books by query params error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if query.Sort == "" {
		query.Sort = "title"
	}
	if query.Order == "" {
		query.Sort = "asc"
	}

	books, err := c.catalogService.ListBooksByCategory(categoryName, query.Page, query.Count, query.Sort, query.Order)
	if err != nil {
		zap.S().Errorf("List books by category error: %v\n", err)
		ctx.JSON(err.Code, gin.H{"error": err.Message})
		return
	}

	ctx.JSON(http.StatusOK, books)
}
