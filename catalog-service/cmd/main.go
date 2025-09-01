package main

import (
	"github.com/Yarik7610/library-backend-common/broker"
	"github.com/Yarik7610/library-backend-common/sharedconstants"
	"github.com/Yarik7610/library-backend/catalog-service/config"
	"github.com/Yarik7610/library-backend/catalog-service/internal/connect"
	"github.com/Yarik7610/library-backend/catalog-service/internal/controller"
	"github.com/Yarik7610/library-backend/catalog-service/internal/repository"
	"github.com/Yarik7610/library-backend/catalog-service/internal/seed"
	"github.com/Yarik7610/library-backend/catalog-service/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
}

func main() {
	err := config.Init()
	if err != nil {
		zap.S().Fatalf("Config load error: %v\n", err)
	}

	db := connect.DB()
	rdb := connect.Cache()
	bookAddedWriter := broker.NewWriter(sharedconstants.BOOK_ADDED_TOPIC)

	bookRepoCache := repository.NewBookRepositoryCache(rdb)
	bookRepo := repository.NewBookRepository(db)
	pageRepo := repository.NewPageRepository(db)
	authorRepo := repository.NewAuthorRepository(db)

	seed.Books(bookRepo, pageRepo, authorRepo)

	catalogService := service.NewCatalogService(db, bookAddedWriter, authorRepo, bookRepoCache, bookRepo, pageRepo)
	catalogController := controller.NewCatalogController(catalogService)

	r := gin.Default()
	catalogGroup := r.Group(sharedconstants.CATALOG_ROUTE)
	{
		catalogGroup.GET(sharedconstants.CATEGORIES_ROUTE, catalogController.GetCategories)
		catalogGroup.GET(sharedconstants.CATEGORIES_ROUTE+"/:categoryName"+sharedconstants.BOOKS_ROUTE, catalogController.ListBooksByCategory)
		catalogGroup.GET(sharedconstants.AUTHORS_ROUTE+"/:authorID"+sharedconstants.BOOKS_ROUTE, catalogController.GetBooksByAuthorID)
		catalogGroup.GET(sharedconstants.BOOKS_ROUTE+sharedconstants.PREVIEW_ROUTE+"/:bookID", catalogController.PreviewBook)
		catalogGroup.GET(sharedconstants.BOOKS_ROUTE+"/:bookID", catalogController.GetBookPage)
		catalogGroup.GET(sharedconstants.BOOKS_ROUTE+sharedconstants.SEARCH_ROUTE, catalogController.SearchBooks)
		catalogGroup.GET(sharedconstants.BOOKS_ROUTE+sharedconstants.NEW_ROUTE, catalogController.GetNewBooks)
		catalogGroup.GET(sharedconstants.BOOKS_ROUTE+sharedconstants.POPULAR_ROUTE, catalogController.GetPopularBooks)
		catalogGroup.GET(sharedconstants.BOOKS_ROUTE+sharedconstants.VIEWS_ROUTE+"/:bookID", catalogController.GetBookViewsCount)

		adminGroup := catalogGroup.Group("")
		{
			adminGroup.DELETE(sharedconstants.BOOKS_ROUTE+"/:bookID", catalogController.DeleteBook)
			adminGroup.POST(sharedconstants.BOOKS_ROUTE, catalogController.AddBook)
			adminGroup.DELETE(sharedconstants.AUTHORS_ROUTE+"/:authorID", catalogController.DeleteAuthor)
			adminGroup.POST(sharedconstants.AUTHORS_ROUTE, catalogController.CreateAuthor)
		}

	}

	if err := r.Run(":" + config.Data.ServerPort); err != nil {
		zap.S().Fatalf("Start error on port %s: %v", config.Data.ServerPort, err)
	}
}
