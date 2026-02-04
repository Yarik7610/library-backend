package http

import (
	"net/http"
	"strconv"

	"github.com/Yarik7610/library-backend-common/sharedconstants"
	"github.com/Yarik7610/library-backend/catalog-service/internal/domain"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/service"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/transport/http/dto"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/transport/http/mapper"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/transport/http/query"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/errs"
	httpInfrastructure "github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/transport/http"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CatalogHandler interface {
	GetCategories(ctx *gin.Context)
	PreviewBook(ctx *gin.Context)
	GetBooksByAuthorID(ctx *gin.Context)
	GetBookPage(ctx *gin.Context)
	AddBook(ctx *gin.Context)
	DeleteBook(ctx *gin.Context)
	CreateAuthor(ctx *gin.Context)
	DeleteAuthor(ctx *gin.Context)
	GetNewBooks(ctx *gin.Context)
	GetBookViewsCount(ctx *gin.Context)
	GetPopularBooks(ctx *gin.Context)
	ListBooksByCategory(ctx *gin.Context)
	SearchBooks(ctx *gin.Context)
}

type catalogHandler struct {
	catalogService service.CatalogService
}

func NewCatalogHandler(catalogService service.CatalogService) CatalogHandler {
	return &catalogHandler{catalogService: catalogService}
}

// GetCategories godoc
//
//	@Summary		Get all book categories
//	@Description	Returns a list of all available book categories
//	@Tags			catalog
//	@Produce		json
//	@Success		200	{array}		string
//	@Failure		500	{object} map[string]string "Internal server error"
//	@Router			/catalog/books/categories [get]
func (c *catalogHandler) GetCategories(ctx *gin.Context) {
	categories, err := c.catalogService.GetCategories()
	if err != nil {
		zap.S().Errorf("List categories error: %v\n", err)
		httpInfrastructure.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, categories)
}

// PreviewBook godoc
//
//	@Summary		Preview a book
//	@Description	Returns preview information for a book
//	@Tags			catalog
//	@Param			bookID	path	uint	true	"Book ID"
//	@Produce		json
//	@Success		200	{object}	dto.Book
//	@Failure		400 {object} map[string]string "Bad request"
//	@Failure		500	{object} map[string]string "Internal server error"
//	@Router			/catalog/books/preview/{bookID} [get]
func (c *catalogHandler) PreviewBook(ctx *gin.Context) {
	bookIDString := ctx.Param("bookID")
	bookID, err := strconv.ParseUint(bookIDString, 10, 64)
	if err != nil {
		zap.S().Errorf("Preview book ID param error: %v\n", err)
		httpInfrastructure.RenderError(ctx, errs.NewBadRequestError(err.Error()))
		return
	}

	userIDString := ctx.GetHeader(sharedconstants.HEADER_USER_ID)
	userID, err := strconv.ParseUint(userIDString, 10, 64)
	if err != nil {
		userID = 0
	}

	bookDomain, err := c.catalogService.PreviewBook(uint(bookID), uint(userID))
	if err != nil {
		zap.S().Errorf("Preview book error: %v\n", err)
		httpInfrastructure.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, mapper.BookDomainToDTO(bookDomain))
}

// GetBooksByAuthorID godoc
//
//	@Summary		Get books by author ID
//	@Description	Returns all books for the given author
//	@Tags			catalog
//	@Param			authorID	path	uint	true	"Author ID"
//	@Produce		json
//	@Success		200	{array}		dto.Book
//	@Failure		400 {object} map[string]string "Bad request"
//	@Failure		500	{object} map[string]string "Internal server error"
//	@Router			/catalog/authors/{authorID}/books [get]
func (c *catalogHandler) GetBooksByAuthorID(ctx *gin.Context) {
	authorIDString := ctx.Param("authorID")
	authorID, err := strconv.ParseUint(authorIDString, 10, 64)
	if err != nil {
		zap.S().Errorf("Get books by author ID param error: %v\n", err)
		httpInfrastructure.RenderError(ctx, errs.NewBadRequestError(err.Error()))
		return
	}

	bookDomains, err := c.catalogService.GetBooksByAuthorID(uint(authorID))
	if err != nil {
		zap.S().Errorf("Get books by author ID error: %v\n", err)
		httpInfrastructure.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, mapper.BookDomainsToDTOs(bookDomains))
}

// GetBookPage godoc
//
//	@Summary		Get a book page
//	@Description	Returns content of a specific book page
//	@Tags			catalog
//	@Param			bookID		path	uint	true	"Book ID"
//	@Param			page	query	int		true	"Page number"
//	@Produce		json
//	@Success		200	{object}	dto.Page
//	@Failure		400 {object} map[string]string "Bad request"
//	@Failure		500	{object} map[string]string "Internal server error"
//	@Router			/catalog/books/{bookID} [get]
func (c *catalogHandler) GetBookPage(ctx *gin.Context) {
	bookIDString := ctx.Param("bookID")
	bookID, err := strconv.ParseUint(bookIDString, 10, 64)
	if err != nil {
		zap.S().Errorf("Get book page book ID param error: %v\n", err)
		httpInfrastructure.RenderError(ctx, errs.NewBadRequestError(err.Error()))
		return
	}

	var query query.GetBookPage
	if err := ctx.ShouldBindQuery(&query); err != nil {
		httpInfrastructure.RenderError(ctx, errs.NewBadRequestError(err.Error()))
		return
	}

	pageDomain, err := c.catalogService.GetBookPage(uint(bookID), query.PageNumber)
	if err != nil {
		zap.S().Errorf("Get book page error: %v\n", err)
		httpInfrastructure.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, mapper.PageDomainToDTO(pageDomain))
}

// AddBook godoc
//
//	@Summary		Add a new book
//	@Description	Creates a new book entry
//	@Tags			catalog
//	@Param			book	body	dto.AddBookRequest	true	"Book info"
//	@Produce		json
//	@Security		BearerAuth
//	@Success		201	{object}	dto.Book
//	@Failure		400 {object} map[string]string "Bad request"
//	@Failure		409 {object} map[string]string "Entity already exists"
//	@Failure		500	{object} map[string]string "Internal server error"
//	@Router			/catalog/books [post]
func (c *catalogHandler) AddBook(ctx *gin.Context) {
	var createBookDTO dto.AddBookRequest
	if err := ctx.ShouldBindJSON(&createBookDTO); err != nil {
		httpInfrastructure.RenderError(ctx, errs.NewBadRequestError(err.Error()))
		return
	}

	bookDomain := mapper.AddBookRequestToDomain(&createBookDTO)
	if err := c.catalogService.AddBook(&bookDomain); err != nil {
		zap.S().Errorf("Add book error: %v\n", err.Error())
		httpInfrastructure.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, mapper.BookDomainToDTO(&bookDomain))
}

// DeleteBook godoc
//
//	@Summary		Delete a book
//	@Description	Deletes a book by ID
//	@Tags			catalog
//	@Param			bookID	path	uint	true	"Book ID"
//	@Produce		json
//	@Security		BearerAuth
//	@Success		204	"No content"
//	@Failure		400 {object} map[string]string "Bad request"
//	@Failure		500	{object} map[string]string "Internal server error"
//	@Router			/catalog/books/{bookID} [delete]
func (c *catalogHandler) DeleteBook(ctx *gin.Context) {
	bookIDString := ctx.Param("bookID")
	bookID, err := strconv.ParseUint(bookIDString, 10, 64)
	if err != nil {
		zap.S().Errorf("Delete book ID param error: %v\n", err)
		httpInfrastructure.RenderError(ctx, errs.NewBadRequestError(err.Error()))
		return
	}

	customErr := c.catalogService.DeleteBook(uint(bookID))
	if customErr != nil {
		zap.S().Errorf("Delete book error: %v\n", err)
		httpInfrastructure.RenderError(ctx, err)
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
//	@Param			author	body	dto.CreateAuthorRequest	true	"Author info"
//	@Produce		json
//	@Security		BearerAuth
//	@Success		201	{object}	dto.Author
//	@Failure		400 {object} map[string]string "Bad request"
//	@Failure		409 {object} map[string]string "Entity already exists"
//	@Failure		500	{object} map[string]string "Internal server error"
//	@Router			/catalog/authors [post]
func (c *catalogHandler) CreateAuthor(ctx *gin.Context) {
	var createAuthorRequestDTO dto.CreateAuthorRequest
	if err := ctx.ShouldBindJSON(&createAuthorRequestDTO); err != nil {
		httpInfrastructure.RenderError(ctx, errs.NewBadRequestError(err.Error()))
		return
	}

	authorDomain := mapper.CreateAuthorRequestDTOToDomain(&createAuthorRequestDTO)
	if err := c.catalogService.CreateAuthor(&authorDomain); err != nil {
		zap.S().Errorf("Create author error: %v\n", err.Error())
		httpInfrastructure.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, mapper.AuthorDomainToDTO(&authorDomain))
}

// DeleteAuthor godoc
//
//	@Summary		Delete an author
//	@Description	Deletes an author by ID
//	@Tags			catalog
//	@Param			authorID	path	uint	true	"Author ID"
//	@Security		BearerAuth
//	@Produce		json
//	@Success		204	"No content"
//	@Failure		400 {object} map[string]string "Bad request"
//	@Failure		500	{object} map[string]string "Internal server error"
//	@Router			/catalog/authors/{authorID} [delete]
func (c *catalogHandler) DeleteAuthor(ctx *gin.Context) {
	authorIDString := ctx.Param("authorID")
	authorID, err := strconv.ParseUint(authorIDString, 10, 64)
	if err != nil {
		httpInfrastructure.RenderError(ctx, errs.NewBadRequestError(err.Error()))
		return
	}

	if err := c.catalogService.DeleteAuthor(uint(authorID)); err != nil {
		zap.S().Errorf("Delete author error: %v\n", err)
		httpInfrastructure.RenderError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
	ctx.Abort()
}

// GetNewBooks godoc
//
//	@Summary		Get new books
//	@Description	Returns list of recently added books
//	@Tags			catalog
//	@Produce		json
//	@Success		200	{array}		dto.Book
//	@Failure		500	{object} map[string]string "Internal server error"
//	@Router			/catalog/books/new [get]
func (c *catalogHandler) GetNewBooks(ctx *gin.Context) {
	newBookDomains, err := c.catalogService.GetNewBooks()
	if err != nil {
		zap.S().Errorf("Get new books error: %v\n", err)
		httpInfrastructure.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, mapper.BookDomainsToDTOs(newBookDomains))
}

// GetBookViewsCount godoc
//
//	@Summary		Get views count for a book
//	@Description	Returns total number of views by different authorized users only
//	@Tags			catalog
//	@Param			bookID	path	uint	true	"Book ID"
//	@Produce		json
//	@Success		200	{object}	map[string]int
//	@Failure		400 {object} map[string]string "Bad request"
//	@Failure		500	{object} map[string]string "Internal server error"
//	@Router			/catalog/books/views/{bookID} [get]
func (c *catalogHandler) GetBookViewsCount(ctx *gin.Context) {
	bookIDString := ctx.Param("bookID")
	bookID, err := strconv.ParseUint(bookIDString, 10, 64)
	if err != nil {
		httpInfrastructure.RenderError(ctx, errs.NewBadRequestError(err.Error()))
		return
	}

	viewsCount, err := c.catalogService.GetBookViewsCount(uint(bookID))
	if err != nil {
		zap.S().Errorf("Get popular books error: %v\n", err)
		httpInfrastructure.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"views": viewsCount})
}

// GetPopularBooks godoc
//
//	@Summary		Get popular books
//	@Description	Returns list of most popular books based on views count by different authorized users only
//	@Tags			catalog
//	@Produce		json
//	@Success		200	{array}		dto.Book
//	@Failure		500	{object} map[string]string "Internal server error"
//	@Router			/catalog/books/popular [get]
func (c *catalogHandler) GetPopularBooks(ctx *gin.Context) {
	popularBookDomains, err := c.catalogService.GetPopularBooks()
	if err != nil {
		zap.S().Errorf("Get popular books error: %v\n", err)
		httpInfrastructure.RenderError(ctx, errs.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, mapper.BookDomainsToDTOs(popularBookDomains))
}

// ListBooksByCategory godoc
//
//	@Summary		List books by category
//	@Description	Returns paginated list of books for the given category
//	@Tags			catalog
//	@Param			categoryName	path	string	true	"Category name"
//	@Param			page	query	int		false	"Page number (min=1, default=1)"
//	@Param			count	query	int		false	"Number of items per page (min=1, max=100, default=20)"
//	@Param			sort			query	string	false	"Sort field (title / year / category, default=title)"
//	@Param			order			query	string	false	"Sort order (asc / desc, default=asc)"
//	@Produce		json
//	@Success		200	{array}		dto.Book
//	@Failure		400 {object} map[string]string "Bad request"
//	@Failure		500	{object} map[string]string "Internal server error"
//	@Router			/catalog/books/categories/{categoryName} [get]
func (c *catalogHandler) ListBooksByCategory(ctx *gin.Context) {
	categoryName := ctx.Param("categoryName")

	var query query.ListBooksByCategory
	if err := ctx.ShouldBindQuery(&query); err != nil {
		httpInfrastructure.RenderError(ctx, errs.NewBadRequestError(err.Error()))
		return
	}

	bookDomains, err := c.catalogService.ListBooksByCategory(categoryName, query.Page, query.Count, query.Sort, query.Order)
	if err != nil {
		zap.S().Errorf("List books by category error: %v\n", err)
		httpInfrastructure.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, mapper.BookDomainsToDTOs(bookDomains))
}

// SearchBooks godoc
//
//	@Summary		Search books
//	@Description	List books by author name and/or title with pagination
//	@Tags			catalog
//	@Param			author	query	string	false	"Author name"
//	@Param			title	query	string	false	"Book title"
//	@Param			page	query	int		false	"Page number (min=1, default=1)"
//	@Param			count	query	int		false	"Number of items per page (min=1, max=100, default=20)"
//	@Param			sort			query	string	false	"Sort field (title / year / category, default=title)"
//	@Param			order			query	string	false	"Sort order (asc / desc, default=asc)"
//	@Produce		json
//	@Success		200	{array}		dto.Book
//	@Failure		400 {object} map[string]string "Bad request"
//	@Failure		500	{object} map[string]string "Internal server error"
//	@Router			/catalog/books/search [get]
func (c *catalogHandler) SearchBooks(ctx *gin.Context) {
	var query query.SearchBooks
	if err := ctx.ShouldBindQuery(&query); err != nil {
		httpInfrastructure.RenderError(ctx, errs.NewBadRequestError(err.Error()))
		return
	}

	if query.Author == "" && query.Title == "" {
		zap.S().Error("Search books error: both author and title are empty")
		httpInfrastructure.RenderError(ctx, errs.NewBadRequestError("Can't have both empty author and title query strings"))
		return
	}

	var bookDomains []domain.Book
	var err error

	if query.Author != "" && query.Title != "" {
		bookDomains, err = c.catalogService.ListBooksByAuthorNameAndTitle(query.Author, query.Title, query.Page, query.Count, query.Sort, query.Order)
	} else if query.Author != "" {
		bookDomains, err = c.catalogService.ListBooksByAuthorName(query.Author, query.Page, query.Count, query.Sort, query.Order)
	} else {
		bookDomains, err = c.catalogService.ListBooksByTitle(query.Title, query.Page, query.Count, query.Sort, query.Order)
	}

	if err != nil {
		zap.S().Errorf("Search books error: %v", err)
		httpInfrastructure.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, mapper.BookDomainsToDTOs(bookDomains))
}
