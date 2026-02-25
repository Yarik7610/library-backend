package http

import (
	"net/http"

	"github.com/Yarik7610/library-backend-common/transport/http/route"
	"github.com/Yarik7610/library-backend/catalog-service/docs"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/config"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func NewRouter(config *config.Config, metricsHandler http.Handler, catalogHandler CatalogHandler) *gin.Engine {
	r := gin.Default()

	r.Use(otelgin.Middleware(config.ServiceName,
		otelgin.WithGinFilter(func(c *gin.Context) bool {
			return c.FullPath() != route.METRICS
		}),
	))

	r.GET(route.METRICS, gin.WrapH(metricsHandler))

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	catalogGroup := r.Group(route.CATALOG)
	{
		bookGroup := catalogGroup.Group(route.BOOKS)
		{
			bookGroup.GET(route.CATEGORIES, catalogHandler.GetBookCategories)
			bookGroup.GET(route.CATEGORIES+"/:categoryName", catalogHandler.ListBooksByCategory)
			bookGroup.GET(route.CATEGORIES+"/exists"+"/:categoryName", catalogHandler.BookCategoryExists) // REMOVE after switch to gRPC
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
