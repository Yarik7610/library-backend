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
		catalogGroup.GET(route.BOOKS+route.CATEGORIES, catalogServiceHandler)
		catalogGroup.GET(route.BOOKS+route.CATEGORIES+"/:categoryName", catalogServiceHandler)
		catalogGroup.GET(route.AUTHORS+"/:authorID"+route.BOOKS, catalogServiceHandler)
		catalogGroup.GET(route.BOOKS+"/:bookID"+route.PREVIEW, core.InjectHeaders(), catalogServiceHandler)
		catalogGroup.GET(route.BOOKS+"/:bookID", catalogServiceHandler)
		catalogGroup.GET(route.BOOKS+route.SEARCH, catalogServiceHandler)
		catalogGroup.GET(route.BOOKS+route.NEW, catalogServiceHandler)
		catalogGroup.GET(route.BOOKS+route.POPULAR, catalogServiceHandler)
		catalogGroup.GET(route.BOOKS+"/:bookID"+route.VIEWS, catalogServiceHandler)

		adminGroup := catalogGroup.Group("")
		adminGroup.Use(middleware.AuthRequired(), middleware.AdminRequired(), core.InjectHeaders())
		{
			adminGroup.DELETE(route.BOOKS+"/:bookID", catalogServiceHandler)
			adminGroup.POST(route.BOOKS, catalogServiceHandler)
			adminGroup.DELETE(route.AUTHORS+"/:authorID", catalogServiceHandler)
			adminGroup.POST(route.AUTHORS, catalogServiceHandler)
		}
	}
}
