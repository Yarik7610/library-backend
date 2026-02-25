package router

import (
	"github.com/Yarik7610/library-backend-common/transport/http/route"
	"github.com/Yarik7610/library-backend/api-gateway/internal/app/middleware"
	"github.com/Yarik7610/library-backend/api-gateway/internal/core"
	"github.com/gin-gonic/gin"
)

func registerCatalogRoutes(r *gin.Engine, catalogMicroserviceHandler gin.HandlerFunc) {
	catalogGroup := r.Group(route.CATALOG)
	{
		bookGroup := catalogGroup.Group(route.BOOKS)
		{
			bookGroup.GET(route.CATEGORIES, catalogMicroserviceHandler)
			bookGroup.GET(route.CATEGORIES+"/:categoryName", catalogMicroserviceHandler)
			bookGroup.GET("/:bookID"+route.PREVIEW, core.InjectHeaders(), catalogMicroserviceHandler)
			bookGroup.GET("/:bookID", catalogMicroserviceHandler)
			bookGroup.GET(route.SEARCH, catalogMicroserviceHandler)
			bookGroup.GET(route.NEW, catalogMicroserviceHandler)
			bookGroup.GET(route.POPULAR, catalogMicroserviceHandler)
			bookGroup.GET("/:bookID"+route.VIEWS, catalogMicroserviceHandler)

			adminGroup := bookGroup.Group("")
			adminGroup.Use(middleware.AuthRequired(), middleware.AdminRequired(), core.InjectHeaders())
			{
				adminGroup.DELETE("/:bookID", catalogMicroserviceHandler)
				adminGroup.POST("", catalogMicroserviceHandler)
			}
		}

		authorGroup := catalogGroup.Group(route.AUTHORS)
		{
			authorGroup.GET("/:authorID"+route.BOOKS, catalogMicroserviceHandler)

			adminGroup := authorGroup.Group("")
			adminGroup.Use(middleware.AuthRequired(), middleware.AdminRequired(), core.InjectHeaders())
			{
				adminGroup.DELETE("/:authorID", catalogMicroserviceHandler)
				adminGroup.POST("", catalogMicroserviceHandler)
			}
		}
	}
}
