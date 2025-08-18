package main

import (
	"github.com/Yarik7610/library-backend-common/sharedconstants"
	"github.com/Yarik7610/library-backend/catalog-service/config"
	"github.com/Yarik7610/library-backend/catalog-service/internal/controller"
	"github.com/Yarik7610/library-backend/catalog-service/internal/model"
	"github.com/Yarik7610/library-backend/catalog-service/internal/repository"
	"github.com/Yarik7610/library-backend/catalog-service/internal/seed"
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

	err = db.AutoMigrate(&model.Author{}, &model.Book{}, &model.Page{})
	if err != nil {
		zap.S().Fatalf("GORM auto migrate error: %v", err)
	}
	zap.S().Info("Successfully made auto migrate")

	bookRepo := repository.NewBookRepository(db)
	pageRepo := repository.NewPageRepository(db)
	authorRepo := repository.NewAuthorRepository(db)

	seed.Books(bookRepo, pageRepo, authorRepo)

	catalogService := service.NewCatalogService(authorRepo, bookRepo, pageRepo)
	catalogController := controller.NewCatalogController(catalogService)

	r := gin.Default()
	catalogRouter := r.Group(sharedconstants.CATALOG_ROUTE)
	{
		catalogRouter.GET(sharedconstants.CATEGORIES_ROUTE, catalogController.GetCategories)
		catalogRouter.GET(sharedconstants.CATEGORIES_ROUTE+"/:categoryName"+sharedconstants.BOOKS_ROUTE, catalogController.ListBooksByCategory)
		catalogRouter.GET(sharedconstants.AUTHORS_ROUTE+"/:authorID"+sharedconstants.BOOKS_ROUTE, catalogController.GetBooksByAuthorID)
		catalogRouter.GET(sharedconstants.BOOKS_ROUTE+sharedconstants.PREVIEW_ROUTE+"/:bookID", catalogController.PreviewBook)
		catalogRouter.GET(sharedconstants.BOOKS_ROUTE+"/:bookID", catalogController.GetBookPage)
		catalogRouter.GET(sharedconstants.BOOKS_ROUTE+sharedconstants.SEARCH_ROUTE, catalogController.SearchBooks)
		// TODO make this for admin only
		catalogRouter.DELETE(sharedconstants.BOOKS_ROUTE+"/:bookID", catalogController.DeleteBook)
		catalogRouter.POST(sharedconstants.BOOKS_ROUTE, catalogController.AddBook)
		catalogRouter.DELETE(sharedconstants.AUTHORS_ROUTE+"/:authorID", catalogController.DeleteAuthor)
		catalogRouter.POST(sharedconstants.AUTHORS_ROUTE, catalogController.CreateAuthor)
	}

	if err := r.Run(":" + config.Data.ServerPort); err != nil {
		zap.S().Fatalf("Start error on port %s: %v", config.Data.ServerPort, err)
	}
}
