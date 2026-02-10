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
		bookGroup := catalogGroup.Group(route.BOOKS)
		{
			bookGroup.GET(route.CATEGORIES, catalogHandler.GetCategories)
			bookGroup.GET(route.CATEGORIES+"/:categoryName", catalogHandler.ListBooksByCategory)
			bookGroup.GET("/:bookID"+route.PREVIEW, catalogHandler.PreviewBook)
			bookGroup.GET("/:bookID", catalogHandler.GetBookPage)
			bookGroup.GET(route.SEARCH, catalogHandler.SearchBooks)
			bookGroup.GET(route.NEW, catalogHandler.GetNewBooks)
			bookGroup.GET(route.POPULAR, catalogHandler.GetPopularBooks)
			bookGroup.GET("/:bookID"+route.VIEWS, catalogHandler.GetBookViewsCount)

			adminGroup := bookGroup.Group("")
			{
				adminGroup.DELETE("/:bookID", catalogHandler.DeleteBook)
				adminGroup.POST("", catalogHandler.AddBook)
			}
		}

		authorGroup := catalogGroup.Group(route.AUTHORS)
		{
			authorGroup.GET("/:authorID"+route.BOOKS, catalogHandler.GetBooksByAuthorID)

			adminGroup := authorGroup.Group("")
			{
				adminGroup.DELETE("/:authorID", catalogHandler.DeleteAuthor)
				adminGroup.POST("", catalogHandler.CreateAuthor)
			}
		}
	}

	return r
}
