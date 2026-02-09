package http

import (
	"github.com/Yarik7610/library-backend-common/transport/http/route"
	"github.com/Yarik7610/library-backend/catalog-service/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(catalogHandler CatalogHandler) *gin.Engine {
	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	catalogGroup := r.Group(route.CATALOG)
	{
		catalogGroup.GET(route.BOOKS+route.CATEGORIES, catalogHandler.GetCategories)
		catalogGroup.GET(route.BOOKS+route.CATEGORIES+"/:categoryName", catalogHandler.ListBooksByCategory)
		catalogGroup.GET(route.AUTHORS+"/:authorID"+route.BOOKS, catalogHandler.GetBooksByAuthorID)
		catalogGroup.GET(route.BOOKS+"/:bookID"+route.PREVIEW, catalogHandler.PreviewBook)
		catalogGroup.GET(route.BOOKS+"/:bookID", catalogHandler.GetBookPage)
		catalogGroup.GET(route.BOOKS+route.SEARCH, catalogHandler.SearchBooks)
		catalogGroup.GET(route.BOOKS+route.NEW, catalogHandler.GetNewBooks)
		catalogGroup.GET(route.BOOKS+route.POPULAR, catalogHandler.GetPopularBooks)
		catalogGroup.GET(route.BOOKS+"/:bookID"+route.VIEWS, catalogHandler.GetBookViewsCount)

		adminGroup := catalogGroup.Group("")
		{
			adminGroup.DELETE(route.BOOKS+"/:bookID", catalogHandler.DeleteBook)
			adminGroup.POST(route.BOOKS, catalogHandler.AddBook)
			adminGroup.DELETE(route.AUTHORS+"/:authorID", catalogHandler.DeleteAuthor)
			adminGroup.POST(route.AUTHORS, catalogHandler.CreateAuthor)
		}
	}

	return r
}
