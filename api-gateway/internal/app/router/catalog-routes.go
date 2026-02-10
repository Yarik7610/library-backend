package router

import (
	"github.com/Yarik7610/library-backend-common/transport/http/route"
	"github.com/Yarik7610/library-backend/api-gateway/internal/app/middleware"
	"github.com/Yarik7610/library-backend/api-gateway/internal/core"
	"github.com/gin-gonic/gin"
)

func registerCatalogRoutes(r *gin.Engine, catalogServiceHandler gin.HandlerFunc) {
	catalogGroup := r.Group(route.CATALOG)
	{
		bookGroup := catalogGroup.Group(route.BOOKS)
		{
			bookGroup.GET(route.CATEGORIES, catalogServiceHandler)
			bookGroup.GET(route.CATEGORIES+"/:categoryName", catalogServiceHandler)
			bookGroup.GET("/:bookID"+route.PREVIEW, core.InjectHeaders(), catalogServiceHandler)
			bookGroup.GET("/:bookID", catalogServiceHandler)
			bookGroup.GET(route.SEARCH, catalogServiceHandler)
			bookGroup.GET(route.NEW, catalogServiceHandler)
			bookGroup.GET(route.POPULAR, catalogServiceHandler)
			bookGroup.GET("/:bookID"+route.VIEWS, catalogServiceHandler)

			adminGroup := bookGroup.Group("")
			adminGroup.Use(middleware.AuthRequired(), middleware.AdminRequired(), core.InjectHeaders())
			{
				adminGroup.DELETE("/:bookID", catalogServiceHandler)
				adminGroup.POST("", catalogServiceHandler)
			}
		}

		authorGroup := catalogGroup.Group(route.AUTHORS)
		{
			authorGroup.GET("/:authorID"+route.BOOKS, catalogServiceHandler)

			adminGroup := authorGroup.Group("")
			adminGroup.Use(middleware.AuthRequired(), middleware.AdminRequired(), core.InjectHeaders())
			{
				adminGroup.DELETE("/:authorID", catalogServiceHandler)
				adminGroup.POST("", catalogServiceHandler)
			}
		}
	}
}
