package http

import (
	"github.com/Yarik7610/library-backend/catalog-service/docs"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/transport/http/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	CATALOG_ROUTE    = "/catalog"
	CATEGORIES_ROUTE = "/categories"
	PREVIEW_ROUTE    = "/preview"
	AUTHORS_ROUTE    = "/authors"
	BOOKS_ROUTE      = "/books"
	SEARCH_ROUTE     = "/search"
	NEW_ROUTE        = "/new"
	POPULAR_ROUTE    = "/popular"
	VIEWS_ROUTE      = "/views"
)

func NewRouter(catalogHandler CatalogHandler) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.AuthContext())

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	catalogGroup := r.Group(CATALOG_ROUTE)
	{
		catalogGroup.GET(BOOKS_ROUTE+CATEGORIES_ROUTE, catalogHandler.GetCategories)
		catalogGroup.GET(BOOKS_ROUTE+CATEGORIES_ROUTE+"/:categoryName", catalogHandler.ListBooksByCategory)
		catalogGroup.GET(AUTHORS_ROUTE+"/:authorID"+BOOKS_ROUTE, catalogHandler.GetBooksByAuthorID)
		catalogGroup.GET(BOOKS_ROUTE+"/:bookID"+PREVIEW_ROUTE, catalogHandler.PreviewBook)
		catalogGroup.GET(BOOKS_ROUTE+"/:bookID", catalogHandler.GetBookPage)
		catalogGroup.GET(BOOKS_ROUTE+SEARCH_ROUTE, catalogHandler.SearchBooks)
		catalogGroup.GET(BOOKS_ROUTE+NEW_ROUTE, catalogHandler.GetNewBooks)
		catalogGroup.GET(BOOKS_ROUTE+POPULAR_ROUTE, catalogHandler.GetPopularBooks)
		catalogGroup.GET(BOOKS_ROUTE+"/:bookID"+VIEWS_ROUTE, catalogHandler.GetBookViewsCount)

		adminGroup := catalogGroup.Group("")
		adminGroup.Use(middleware.AdminRequired())
		{
			adminGroup.DELETE(BOOKS_ROUTE+"/:bookID", catalogHandler.DeleteBook)
			adminGroup.POST(BOOKS_ROUTE, catalogHandler.AddBook)
			adminGroup.DELETE(AUTHORS_ROUTE+"/:authorID", catalogHandler.DeleteAuthor)
			adminGroup.POST(AUTHORS_ROUTE, catalogHandler.CreateAuthor)
		}
	}

	return r
}
