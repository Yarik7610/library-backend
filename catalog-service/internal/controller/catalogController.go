package controller

import (
	"net/http"
	"strconv"

	"github.com/Yarik7610/library-backend-common/custom"
	"github.com/Yarik7610/library-backend/catalog-service/internal/dto"
	"github.com/Yarik7610/library-backend/catalog-service/internal/query"
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
	GetBookPage(ctx *gin.Context)
	DeleteBook(ctx *gin.Context)
	AddBook(ctx *gin.Context)
	DeleteAuthor(ctx *gin.Context)
	CreateAuthor(ctx *gin.Context)
}

type catalogController struct {
	catalogService service.CatalogService
}

func NewCatalogController(catalogService service.CatalogService) CatalogController {
	return &catalogController{catalogService: catalogService}
}

func (c *catalogController) PreviewBook(ctx *gin.Context) {
	bookIDString := ctx.Param("bookID")
	bookID, err := strconv.ParseUint(bookIDString, 10, 64)
	if err != nil {
		zap.S().Errorf("Preview book ID param error: %v\n", err)
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
	authorID, err := strconv.ParseUint(authorIDString, 10, 64)
	if err != nil {
		zap.S().Errorf("Get books by author ID param error: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var customErr *custom.Err
	books, customErr := c.catalogService.GetBooksByAuthorID(uint(authorID))
	if customErr != nil {
		zap.S().Errorf("Get books by author ID error: %v\n", err)
		ctx.JSON(customErr.Code, gin.H{"error": customErr.Message})
		return
	}

	ctx.JSON(http.StatusOK, books)
}

func (c *catalogController) SearchBooks(ctx *gin.Context) {
	var q query.SearchBooks
	if err := ctx.ShouldBindQuery(&q); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if q.Author == "" && q.Title == "" {
		zap.S().Error("Search books error: both author and title are empty")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Can't have both empty author and title query strings"})
		return
	}

	sort, order := initOrderParams(q.Sort, q.Order)

	var books []dto.ListedBooks
	var err *custom.Err

	if q.Author != "" && q.Title != "" {
		books, err = c.catalogService.ListBooksByAuthorNameAndTitle(q.Author, q.Title, q.Page, q.Count, sort, order)
	} else if q.Author != "" {
		books, err = c.catalogService.ListBooksByAuthorName(q.Author, q.Page, q.Count, sort, order)
	} else {
		books, err = c.catalogService.ListBooksByTitle(q.Title, q.Page, q.Count, sort, order)
	}

	if err != nil {
		zap.S().Errorf("Search books error: %v", err.Message)
		ctx.JSON(err.Code, gin.H{"error": err.Message})
		return
	}

	ctx.JSON(http.StatusOK, books)
}

func (c *catalogController) ListBooksByCategory(ctx *gin.Context) {
	categoryName := ctx.Param("categoryName")

	var q query.ListBooksByCategory
	if err := ctx.ShouldBindQuery(&q); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sort, order := initOrderParams(q.Sort, q.Order)

	books, err := c.catalogService.ListBooksByCategory(categoryName, q.Page, q.Count, sort, order)
	if err != nil {
		zap.S().Errorf("List books by category error: %v\n", err)
		ctx.JSON(err.Code, gin.H{"error": err.Message})
		return
	}

	ctx.JSON(http.StatusOK, books)
}

func (c *catalogController) GetBookPage(ctx *gin.Context) {
	bookIDString := ctx.Param("bookID")
	bookID, err := strconv.ParseUint(bookIDString, 10, 64)
	if err != nil {
		zap.S().Errorf("Delete book ID param error: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var q query.GetBookPage
	if err := ctx.ShouldBindQuery(&q); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var customErr *custom.Err
	page, customErr := c.catalogService.GetBookPage(uint(bookID), q.PageNumber)
	if customErr != nil {
		zap.S().Errorf("Get book page error: %v\n", err)
		ctx.JSON(customErr.Code, gin.H{"error": customErr.Message})
		return
	}

	ctx.JSON(http.StatusOK, page)
}

func (c *catalogController) DeleteBook(ctx *gin.Context) {
	bookIDString := ctx.Param("bookID")
	bookID, err := strconv.ParseUint(bookIDString, 10, 64)
	if err != nil {
		zap.S().Errorf("Delete book ID param error: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customErr := c.catalogService.DeleteBook(uint(bookID))
	if customErr != nil {
		zap.S().Errorf("Delete book error: %v\n", err)
		ctx.JSON(customErr.Code, gin.H{"error": customErr.Message})
		return
	}

	ctx.Status(http.StatusNoContent)
	ctx.Abort()
}

func (c *catalogController) AddBook(ctx *gin.Context) {
	var createBookDTO dto.AddBook
	if err := ctx.ShouldBindJSON(&createBookDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, customErr := c.catalogService.AddBook(&createBookDTO)
	if customErr != nil {
		zap.S().Errorf("Create book error: %v\n", customErr.Error())
		ctx.JSON(customErr.Code, gin.H{"error": customErr.Message})
		return
	}

	ctx.JSON(http.StatusCreated, book)
}

func (c *catalogController) DeleteAuthor(ctx *gin.Context) {
	authorIDString := ctx.Param("authorID")
	authorID, err := strconv.ParseUint(authorIDString, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customErr := c.catalogService.DeleteAuthor(uint(authorID))
	if customErr != nil {
		zap.S().Errorf("Delete author error: %v\n", err)
		ctx.JSON(customErr.Code, gin.H{"error": customErr.Message})
		return
	}

	ctx.Status(http.StatusNoContent)
	ctx.Abort()
}

func (c *catalogController) CreateAuthor(ctx *gin.Context) {
	var createAuthorDTO dto.CreateAuthor
	if err := ctx.ShouldBindJSON(&createAuthorDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	author, customErr := c.catalogService.CreateAuthor(createAuthorDTO.Fullname)
	if customErr != nil {
		zap.S().Errorf("Create author error: %v\n", customErr.Error())
		ctx.JSON(customErr.Code, gin.H{"error": customErr.Message})
		return
	}

	ctx.JSON(http.StatusCreated, author)
}

func initOrderParams(sort, order string) (string, string) {
	if sort == "" {
		sort = "title"
	}
	if order == "" {
		order = "asc"
	}
	return sort, order
}
