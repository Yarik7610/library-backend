package http

import (
	"net/http"
	"strconv"

	"github.com/Yarik7610/library-backend/catalog-service/internal/domain"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/service"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/transport/http/dto"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/transport/http/mapper"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/transport/http/query"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/errs"
	httpInfrastructure "github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/transport/http"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/transport/http/header"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CatalogHandler interface {
	GetCategories(c *gin.Context)
	PreviewBook(c *gin.Context)
	GetBooksByAuthorID(c *gin.Context)
	GetBookPage(c *gin.Context)
	AddBook(c *gin.Context)
	DeleteBook(c *gin.Context)
	CreateAuthor(c *gin.Context)
	DeleteAuthor(c *gin.Context)
	GetNewBooks(c *gin.Context)
	GetBookViewsCount(c *gin.Context)
	GetPopularBooks(c *gin.Context)
	ListBooksByCategory(c *gin.Context)
	SearchBooks(c *gin.Context)
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
//	@Failure		500	{object} 	dto.Error "Internal server error"
//	@Router			/catalog/books/categories [get]
func (h *catalogHandler) GetCategories(c *gin.Context) {
	ctx := c.Request.Context()

	categories, err := h.catalogService.GetCategories(ctx)
	if err != nil {
		zap.S().Errorf("Get categories error: %v\n", err)
		httpInfrastructure.RenderError(c, err)
		return
	}

	c.JSON(http.StatusOK, categories)
}

// PreviewBook godoc
//
//	@Summary		Preview a book
//	@Description	Returns preview information for a book
//	@Tags			catalog
//	@Param			bookID	path	uint	true	"Book ID"
//	@Produce		json
//	@Success		200	{object}	dto.Book
//	@Failure		400 {object} 	dto.Error "Bad request"
//	@Failure		404 {object} 	dto.Error "Entity not found"
//	@Failure		500	{object} 	dto.Error "Internal server error"
//	@Router			/catalog/books/{bookID}/preview [get]
func (h *catalogHandler) PreviewBook(c *gin.Context) {
	ctx := c.Request.Context()

	bookIDString := c.Param("bookID")
	bookID, err := strconv.ParseUint(bookIDString, 10, 64)
	if err != nil {
		httpInfrastructure.RenderError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	userID, err := header.GetUserID(c)
	if err != nil {
		userID = 0
	}

	bookDomain, err := h.catalogService.PreviewBook(ctx, uint(bookID), uint(userID))
	if err != nil {
		zap.S().Errorf("Preview book error: %v\n", err)
		httpInfrastructure.RenderError(c, err)
		return
	}

	c.JSON(http.StatusOK, mapper.BookDomainToDTO(bookDomain))
}

// GetBooksByAuthorID godoc
//
//	@Summary		Get books by author ID
//	@Description	Returns all books for the given author
//	@Tags			catalog
//	@Param			authorID	path	uint	true	"Author ID"
//	@Produce		json
//	@Success		200	{array}		dto.Book
//	@Failure		400 {object} 	dto.Error "Bad request"
//	@Failure		500	{object} 	dto.Error "Internal server error"
//	@Router			/catalog/authors/{authorID}/books [get]
func (h *catalogHandler) GetBooksByAuthorID(c *gin.Context) {
	ctx := c.Request.Context()

	authorIDString := c.Param("authorID")
	authorID, err := strconv.ParseUint(authorIDString, 10, 64)
	if err != nil {
		httpInfrastructure.RenderError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	bookDomains, err := h.catalogService.GetBooksByAuthorID(ctx, uint(authorID))
	if err != nil {
		zap.S().Errorf("Get books by author ID error: %v\n", err)
		httpInfrastructure.RenderError(c, err)
		return
	}

	c.JSON(http.StatusOK, mapper.BookDomainsToDTOs(bookDomains))
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
//	@Failure		400 {object} 	dto.Error "Bad request"
//	@Failure		404 {object} 	dto.Error "Entity not found"
//	@Failure		500	{object}	dto.Error "Internal server error"
//	@Router			/catalog/books/{bookID} [get]
func (h *catalogHandler) GetBookPage(c *gin.Context) {
	ctx := c.Request.Context()

	bookIDString := c.Param("bookID")
	bookID, err := strconv.ParseUint(bookIDString, 10, 64)
	if err != nil {
		httpInfrastructure.RenderError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	var query query.GetBookPage
	if err := c.ShouldBindQuery(&query); err != nil {
		httpInfrastructure.RenderError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	pageDomain, err := h.catalogService.GetBookPage(ctx, uint(bookID), query.PageNumber)
	if err != nil {
		zap.S().Errorf("Get book page error: %v\n", err)
		httpInfrastructure.RenderError(c, err)
		return
	}

	c.JSON(http.StatusOK, mapper.PageDomainToDTO(pageDomain))
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
//	@Failure		400 {object} 	dto.Error "Bad request"
//	@Failure		401 {object} 	dto.Error "The token is missing, invalid or expired"
//	@Failure		403 {object} 	dto.Error "The token is valid, but lacks permission"
//	@Failure		404 {object} 	dto.Error "Entity not found"
//	@Failure		409 {object} 	dto.Error "Entity already exists"
//	@Failure		500	{object} 	dto.Error "Internal server error"
//	@Router			/catalog/books [post]
func (h *catalogHandler) AddBook(c *gin.Context) {
	ctx := c.Request.Context()

	var createBookDTO dto.AddBookRequest
	if err := c.ShouldBindJSON(&createBookDTO); err != nil {
		httpInfrastructure.RenderError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	bookDomain := mapper.AddBookRequestToDomain(&createBookDTO)
	if err := h.catalogService.AddBook(ctx, &bookDomain); err != nil {
		zap.S().Errorf("Add book error: %v\n", err.Error())
		httpInfrastructure.RenderError(c, err)
		return
	}

	c.JSON(http.StatusCreated, mapper.BookDomainToDTO(&bookDomain))
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
//	@Failure		400 {object}	dto.Error "Bad request"
//	@Failure		401 {object}	dto.Error "The token is missing, invalid or expired"
//	@Failure		403 {object}	dto.Error "The token is valid, but lacks permission"
//	@Failure		500	{object}	dto.Error "Internal server error"
//	@Router			/catalog/books/{bookID} [delete]
func (h *catalogHandler) DeleteBook(c *gin.Context) {
	ctx := c.Request.Context()

	bookIDString := c.Param("bookID")
	bookID, err := strconv.ParseUint(bookIDString, 10, 64)
	if err != nil {
		httpInfrastructure.RenderError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	err = h.catalogService.DeleteBook(ctx, uint(bookID))
	if err != nil {
		zap.S().Errorf("Delete book error: %v\n", err)
		httpInfrastructure.RenderError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
	c.Abort()
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
//	@Failure		400 {object}	dto.Error "Bad request"
//	@Failure		401 {object}	dto.Error "The token is missing, invalid or expired"
//	@Failure		403 {object}	dto.Error "The token is valid, but lacks permission"
//	@Failure		409 {object}	dto.Error "Entity already exists"
//	@Failure		500	{object}	dto.Error "Internal server error"
//	@Router			/catalog/authors [post]
func (h *catalogHandler) CreateAuthor(c *gin.Context) {
	ctx := c.Request.Context()

	var createAuthorRequestDTO dto.CreateAuthorRequest
	if err := c.ShouldBindJSON(&createAuthorRequestDTO); err != nil {
		httpInfrastructure.RenderError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	authorDomain := mapper.CreateAuthorRequestDTOToDomain(&createAuthorRequestDTO)
	if err := h.catalogService.CreateAuthor(ctx, &authorDomain); err != nil {
		zap.S().Errorf("Create author error: %v\n", err.Error())
		httpInfrastructure.RenderError(c, err)
		return
	}

	c.JSON(http.StatusCreated, mapper.AuthorDomainToDTO(&authorDomain))
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
//	@Failure		400 {object}	dto.Error "Bad request"
//	@Failure		401 {object}	dto.Error "The token is missing, invalid or expired"
//	@Failure		403 {object}	dto.Error "The token is valid, but lacks permission"
//	@Failure		500	{object}	dto.Error "Internal server error"
//	@Router			/catalog/authors/{authorID} [delete]
func (h *catalogHandler) DeleteAuthor(c *gin.Context) {
	ctx := c.Request.Context()

	authorIDString := c.Param("authorID")
	authorID, err := strconv.ParseUint(authorIDString, 10, 64)
	if err != nil {
		httpInfrastructure.RenderError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	if err := h.catalogService.DeleteAuthor(ctx, uint(authorID)); err != nil {
		zap.S().Errorf("Delete author error: %v\n", err)
		httpInfrastructure.RenderError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
	c.Abort()
}

// GetNewBooks godoc
//
//	@Summary		Get new books
//	@Description	Returns list of recently added books
//	@Tags			catalog
//	@Produce		json
//	@Success		200	{array}		dto.Book
//	@Failure		500	{object}	dto.Error "Internal server error"
//	@Router			/catalog/books/new [get]
func (h *catalogHandler) GetNewBooks(c *gin.Context) {
	ctx := c.Request.Context()

	newBookDomains, err := h.catalogService.GetNewBooks(ctx)
	if err != nil {
		zap.S().Errorf("Get new books error: %v\n", err)
		httpInfrastructure.RenderError(c, err)
		return
	}

	c.JSON(http.StatusOK, mapper.BookDomainsToDTOs(newBookDomains))
}

// GetBookViewsCount godoc
//
//	@Summary		Get views count for a book
//	@Description	Returns total number of views by different authorized users only. If book doesn't exist it still will return zero views count
//	@Tags			catalog
//	@Param			bookID	path	uint	true	"Book ID"
//	@Produce		json
//	@Success		200	{object}	dto.BookViews
//	@Failure		400 {object}	dto.Error "Bad request"
//	@Failure		500	{object}	dto.Error "Internal server error"
//	@Router			/catalog/books/{bookID}/views [get]
func (h *catalogHandler) GetBookViewsCount(c *gin.Context) {
	ctx := c.Request.Context()

	bookIDString := c.Param("bookID")
	bookID, err := strconv.ParseUint(bookIDString, 10, 64)
	if err != nil {
		httpInfrastructure.RenderError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	viewsCount, err := h.catalogService.GetBookViewsCount(ctx, uint(bookID))
	if err != nil {
		zap.S().Errorf("Get book views count error: %v\n", err)
		httpInfrastructure.RenderError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.BookViews{Views: viewsCount})
}

// GetPopularBooks godoc
//
//	@Summary		Get popular books
//	@Description	Returns list of most popular books based on views count by different authorized users only
//	@Tags			catalog
//	@Produce		json
//	@Success		200	{array}		dto.Book
//	@Failure		500	{object}	dto.Error "Internal server error"
//	@Router			/catalog/books/popular [get]
func (h *catalogHandler) GetPopularBooks(c *gin.Context) {
	ctx := c.Request.Context()

	popularBookDomains, err := h.catalogService.GetPopularBooks(ctx)
	if err != nil {
		zap.S().Errorf("Get popular books error: %v\n", err)
		httpInfrastructure.RenderError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, mapper.BookDomainsToDTOs(popularBookDomains))
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
//	@Failure		400 {object}	dto.Error "Bad request"
//	@Failure		500	{object}	dto.Error "Internal server error"
//	@Router			/catalog/books/categories/{categoryName} [get]
func (h *catalogHandler) ListBooksByCategory(c *gin.Context) {
	ctx := c.Request.Context()

	categoryName := c.Param("categoryName")

	var query query.ListBooksByCategory
	if err := c.ShouldBindQuery(&query); err != nil {
		httpInfrastructure.RenderError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	bookDomains, err := h.catalogService.ListBooksByCategory(ctx, categoryName, query.Page, query.Count, query.Sort, query.Order)
	if err != nil {
		zap.S().Errorf("List books by category error: %v\n", err)
		httpInfrastructure.RenderError(c, err)
		return
	}

	c.JSON(http.StatusOK, mapper.BookDomainsToDTOs(bookDomains))
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
//	@Failure		400 {object}	dto.Error "Bad request"
//	@Failure		500	{object}	dto.Error "Internal server error"
//	@Router			/catalog/books/search [get]
func (h *catalogHandler) SearchBooks(c *gin.Context) {
	ctx := c.Request.Context()

	var query query.SearchBooks
	if err := c.ShouldBindQuery(&query); err != nil {
		httpInfrastructure.RenderError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	if query.Author == "" && query.Title == "" {
		httpInfrastructure.RenderError(c, errs.NewBadRequestError("Can't have both empty author and title query strings"))
		return
	}

	var bookDomains []domain.Book
	var err error

	if query.Author != "" && query.Title != "" {
		bookDomains, err = h.catalogService.ListBooksByAuthorNameAndTitle(ctx, query.Author, query.Title, query.Page, query.Count, query.Sort, query.Order)
	} else if query.Author != "" {
		bookDomains, err = h.catalogService.ListBooksByAuthorName(ctx, query.Author, query.Page, query.Count, query.Sort, query.Order)
	} else {
		bookDomains, err = h.catalogService.ListBooksByTitle(ctx, query.Title, query.Page, query.Count, query.Sort, query.Order)
	}

	if err != nil {
		zap.S().Errorf("Search books error: %v", err)
		httpInfrastructure.RenderError(c, err)
		return
	}

	c.JSON(http.StatusOK, mapper.BookDomainsToDTOs(bookDomains))
}
