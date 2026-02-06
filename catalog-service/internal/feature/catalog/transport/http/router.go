package http

import (
	"github.com/Yarik7610/library-backend-common/sharedconstants"
	"github.com/Yarik7610/library-backend/catalog-service/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(catalogHandler CatalogHandler) *gin.Engine {
	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	catalogGroup := r.Group(sharedconstants.CATALOG_ROUTE)
	{
		catalogGroup.GET(sharedconstants.BOOKS_ROUTE+sharedconstants.CATEGORIES_ROUTE, catalogHandler.GetCategories)
		catalogGroup.GET(sharedconstants.BOOKS_ROUTE+sharedconstants.CATEGORIES_ROUTE+"/:categoryName", catalogHandler.ListBooksByCategory)
		catalogGroup.GET(sharedconstants.AUTHORS_ROUTE+"/:authorID"+sharedconstants.BOOKS_ROUTE, catalogHandler.GetBooksByAuthorID)
		catalogGroup.GET(sharedconstants.BOOKS_ROUTE+"/:bookID"+sharedconstants.PREVIEW_ROUTE, catalogHandler.PreviewBook)
		catalogGroup.GET(sharedconstants.BOOKS_ROUTE+"/:bookID", catalogHandler.GetBookPage)
		catalogGroup.GET(sharedconstants.BOOKS_ROUTE+sharedconstants.SEARCH_ROUTE, catalogHandler.SearchBooks)
		catalogGroup.GET(sharedconstants.BOOKS_ROUTE+sharedconstants.NEW_ROUTE, catalogHandler.GetNewBooks)
		catalogGroup.GET(sharedconstants.BOOKS_ROUTE+sharedconstants.POPULAR_ROUTE, catalogHandler.GetPopularBooks)
		catalogGroup.GET(sharedconstants.BOOKS_ROUTE+"/:bookID"+sharedconstants.VIEWS_ROUTE, catalogHandler.GetBookViewsCount)

		adminGroup := catalogGroup.Group("")
		{
			adminGroup.DELETE(sharedconstants.BOOKS_ROUTE+"/:bookID", catalogHandler.DeleteBook)
			adminGroup.POST(sharedconstants.BOOKS_ROUTE, catalogHandler.AddBook)
			adminGroup.DELETE(sharedconstants.AUTHORS_ROUTE+"/:authorID", catalogHandler.DeleteAuthor)
			adminGroup.POST(sharedconstants.AUTHORS_ROUTE, catalogHandler.CreateAuthor)
		}
	}

	return r
}
