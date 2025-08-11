package controller

import (
	"github.com/Yarik7610/library-backend/catalog-service/internal/service"
	"github.com/gin-gonic/gin"
)

type CatalogController interface {
	PreviewBook(ctx *gin.Context)
	ListCategories(ctx *gin.Context)
}

type catalogController struct {
	catalogService service.CatalogService
}

func NewCatalogController(catalogService service.CatalogService) CatalogController {
	return &catalogController{catalogService: catalogService}
}

func (c *catalogController) PreviewBook(ctx *gin.Context) {

}

func (c *catalogController) ListCategories(ctx *gin.Context) {

}
