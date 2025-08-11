package controller

import (
	"github.com/Yarik7610/library-backend/catalog-service/internal/service"
	"github.com/gin-gonic/gin"
)

type CategoryController interface {
	ListCategories(ctx *gin.Context)
}

type categoryController struct {
	categoryService service.CategoryService
}

func NewCategoryController(categoryService service.CategoryService) CategoryController {
	return &categoryController{categoryService: categoryService}
}

func (c *categoryController) ListCategories(ctx *gin.Context) {

}
