package main

import (
	"github.com/Yarik7610/library-backend-common/sharedconstants"
	"github.com/Yarik7610/library-backend/catalog-service/config"
	"github.com/Yarik7610/library-backend/catalog-service/internal/controller"
	"github.com/Yarik7610/library-backend/catalog-service/internal/model"
	"github.com/Yarik7610/library-backend/catalog-service/internal/repository"
	"github.com/Yarik7610/library-backend/catalog-service/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
}

func main() {
	err := config.Init()
	if err != nil {
		zap.S().Fatalf("Config load error: %v\n", err)
	}

	db, err := gorm.Open(postgres.Open(config.Data.PostgresURL), &gorm.Config{})
	if err != nil {
		zap.S().Fatalf("GORM open error: %v\n", err)
	}
	zap.S().Info("Successfully connected to Postgres")

	err = db.AutoMigrate(&model.Book{}, &model.Category{}, &model.Page{})
	if err != nil {
		zap.S().Fatalf("GORM auto migrate error: %v", err)
	}
	zap.S().Info("Successfully made auto migrate")

	bookRepo := repository.NewBookRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	bookService := service.NewBookService(bookRepo)
	categoryService := service.NewCategoryService(categoryRepo)

	bookController := controller.NewBookController(bookService)
	categoryController := controller.NewCategoryController(categoryService)

	r := gin.Default()
	catalogRouter := r.Group(sharedconstants.CATALOG_ROUTE)
	{
		catalogRouter.GET(sharedconstants.CATEGORIES_ROUTE, categoryController.ListCategories)
		catalogRouter.GET(sharedconstants.PREVIEW_ROUTE, bookController.PreviewBook)
	}

	if err := r.Run(":" + config.Data.ServerPort); err != nil {
		zap.S().Fatalf("User-service start error on port %s: %v", config.Data.ServerPort, err)
	}
}
