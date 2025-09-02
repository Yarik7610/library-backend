package controller

import (
	"net/http"
	"strconv"

	"github.com/Yarik7610/library-backend-common/custom"
	"github.com/Yarik7610/library-backend-common/sharedconstants"
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
	GetNewBooks(ctx *gin.Context)
	GetPopularBooks(ctx *gin.Context)
	GetBookViewsCount(ctx *gin.Context)
}

type catalogController struct {
	catalogService service.CatalogService
}

func NewCatalogController(catalogService service.CatalogService) CatalogController {
	return &catalogController{catalogService: catalogService}
}

// PreviewBook godoc
//
//	@Summary		Preview a book
//	@Description	Returns preview information for a book
//	@Tags			catalog
//	@Param			bookID	path	uint	true	"Book ID"
//	@Produce		json
//	@Success		200	{object}	model.Book
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/catalog/books/preview/{bookID} [get]
func (c *catalogController) PreviewBook(ctx *gin.Context) {
	bookIDString := ctx.Param("bookID")
	bookID, err := strconv.ParseUint(bookIDString, 10, 64)
	if err != nil {
		zap.S().Errorf("Preview book ID param error: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDString := ctx.GetHeader(sharedconstants.HEADER_USER_ID)
	userID, err := strconv.ParseUint(userIDString, 10, 64)
	if err != nil {
		userID = 0
	}

	book, customErr := c.catalogService.PreviewBook(uint(bookID), uint(userID))
	if customErr != nil {
		zap.S().Errorf("Preview book error: %v\n", err)
		ctx.JSON(customErr.Code, gin.H{"error": customErr.Message})
		return
	}

	ctx.JSON(http.StatusOK, book)
}

// GetCategories godoc
//
//	@Summary		Get all book categories
//	@Description	Returns a list of all available book categories
//	@Tags			catalog
//	@Produce		json
//	@Success		200	{array}		string
//	@Failure		500	{object}	map[string]string
//	@Router			/catalog/categories [get]
func (c *catalogController) GetCategories(ctx *gin.Context) {
	categories, err := c.catalogService.GetCategories()
	if err != nil {
		zap.S().Errorf("List categories error: %v\n", err)
		ctx.JSON(err.Code, gin.H{"error": err.Message})
		return
	}

	ctx.JSON(http.StatusOK, categories)
}

// GetBooksByAuthorID godoc
//
//	@Summary		Get books by author ID
//	@Description	Returns all books for the given author
//	@Tags			catalog
//	@Param			authorID	path	uint	true	"Author ID"
//	@Produce		json
//	@Success		200	{array}		dto.ListedBooks
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/catalog/authors/{authorID}/books [get]
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

// SearchBooks godoc
//
//	@Summary		Search books
//	@Description	Search books by author name and/or title with pagination
//	@Tags			catalog
//	@Param			author	query	string	false	"Author name"
//	@Param			title	query	string	false	"Book title"
//	@Param			page	query	int		false	"Page number"
//	@Param			count	query	int		false	"Number of items per page"
//	@Param			sort	query	string	false	"Sort field"
//	@Param			order	query	string	false	"Sort order (asc/desc)"
//	@Produce		json
//	@Success		200	{array}		dto.ListedBooks
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/catalog/books/search [get]
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

// ListBooksByCategory godoc
//
//	@Summary		List books by category
//	@Description	Returns paginated list of books for the given category
//	@Tags			catalog
//	@Param			categoryName	path	string	true	"Category name"
//	@Param			page			query	int		false	"Page number"
//	@Param			count			query	int		false	"Number of items per page"
//	@Param			sort			query	string	false	"Sort field"
//	@Param			order			query	string	false	"Sort order (asc/desc)"
//	@Produce		json
//	@Success		200	{array}		dto.ListedBooks
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/catalog/categories/{categoryName}/books [get]
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

// GetBookPage godoc
//
//	@Summary		Get a book page
//	@Description	Returns content of a specific book page
//	@Tags			catalog
//	@Param			bookID		path	uint	true	"Book ID"
//	@Param			pageNumber	query	int		true	"Page number"
//	@Produce		json
//	@Success		200	{object}	model.Page
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/catalog/books/{bookID} [get]
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

// DeleteBook godoc
//
//	@Summary		Delete a book
//	@Description	Deletes a book by ID
//	@Tags			catalog
//	@Param			bookID	path	uint	true	"Book ID"
//	@Produce		json
//	@Security		BearerAuth
//	@Success		204	""
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/catalog/books/{bookID} [delete]
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

// AddBook godoc
//
//	@Summary		Add a new book
//	@Description	Creates a new book entry
//	@Tags			catalog
//	@Param			book	body	dto.AddBook	true	"Book info"
//	@Produce		json
//	@Security		BearerAuth
//	@Success		201	{object}	model.Book
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/catalog/books [post]
func (c *catalogController) AddBook(ctx *gin.Context) {
	var createBookDTO dto.AddBook
	if err := ctx.ShouldBindJSON(&createBookDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, customErr := c.catalogService.AddBook(&createBookDTO)
	if customErr != nil {
		zap.S().Errorf("Add book error: %v\n", customErr.Error())
		ctx.JSON(customErr.Code, gin.H{"error": customErr.Message})
		return
	}

	ctx.JSON(http.StatusCreated, book)
}

// DeleteAuthor godoc
//
//	@Summary		Delete an author
//	@Description	Deletes an author by ID
//	@Tags			catalog
//	@Param			authorID	path	uint	true	"Author ID"
//	@Security		BearerAuth
//	@Produce		json
//	@Success		204	""
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/catalog/authors/{authorID} [delete]
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

// CreateAuthor godoc
//
//	@Summary		Create a new author
//	@Description	Adds a new author
//	@Tags			catalog
//	@Param			author	body	dto.CreateAuthor	true	"Author info"
//	@Produce		json
//	@Security		BearerAuth
//	@Success		201	{object}	model.Author
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/catalog/authors [post]
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

// GetNewBooks godoc
//
//	@Summary		Get new books
//	@Description	Returns list of recently added books
//	@Tags			catalog
//	@Produce		json
//	@Success		200	{array}		dto.ListedBooks
//	@Failure		500	{object}	map[string]string
//	@Router			/catalog/books/new [get]
func (c *catalogController) GetNewBooks(ctx *gin.Context) {
	newBooks, err := c.catalogService.GetNewBooks()
	if err != nil {
		zap.S().Errorf("Get new books error: %v\n", err)
		ctx.JSON(err.Code, gin.H{"error": err.Message})
		return
	}

	ctx.JSON(http.StatusOK, newBooks)
}

// GetPopularBooks godoc
//
//	@Summary		Get popular books
//	@Description	Returns list of most popular books
//	@Tags			catalog
//	@Produce		json
//	@Success		200	{array}		dto.ListedBooks
//	@Failure		500	{object}	map[string]string
//	@Router			/catalog/books/popular [get]
func (c *catalogController) GetPopularBooks(ctx *gin.Context) {
	popularBooks, err := c.catalogService.GetPopularBooks()
	if err != nil {
		zap.S().Errorf("Get popular books error: %v\n", err)
		ctx.JSON(err.Code, gin.H{"error": err.Message})
		return
	}

	ctx.JSON(http.StatusOK, popularBooks)
}

// GetBookViewsCount godoc
//
//	@Summary		Get views count for a book
//	@Description	Returns total number of views
//	@Tags			catalog
//	@Param			bookID	path	uint	true	"Book ID"
//	@Produce		json
//	@Success		200	{object}	map[string]int
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/catalog/books/views/{bookID} [get]
func (c *catalogController) GetBookViewsCount(ctx *gin.Context) {
	bookIDString := ctx.Param("bookID")
	bookID, err := strconv.ParseUint(bookIDString, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	viewsCount, customErr := c.catalogService.GetBookViewsCount(uint(bookID))
	if customErr != nil {
		zap.S().Errorf("Get popular books error: %v\n", err)
		ctx.JSON(customErr.Code, gin.H{"error": customErr.Message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"views": viewsCount})
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
