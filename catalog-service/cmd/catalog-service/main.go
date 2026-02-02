package main

import (
	"github.com/Yarik7610/library-backend-common/broker"
	"github.com/Yarik7610/library-backend-common/sharedconstants"

	docs "github.com/Yarik7610/library-backend/catalog-service/docs"
	postgresRepositories "github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/postgres"
	redisRepositories "github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/redis"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/service"
	controller "github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/transport/http"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/storage/postgres"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/storage/postgres/seed"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/storage/redis"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	postgresDB := postgres.Connect()
	rdb := redis.Connect()

	bookAddedWriter := broker.NewWriter(sharedconstants.BOOK_ADDED_TOPIC)

	bookRepoCache := redisRepositories.NewBookRepository(rdb)
	bookRepo := postgresRepositories.NewBookRepository(postgresDB)
	pageRepo := postgresRepositories.NewPageRepository(postgresDB)
	authorRepo := postgresRepositories.NewAuthorRepository(postgresDB)

	seed.Books(bookRepo, pageRepo, authorRepo)

	catalogService := service.NewCatalogService(postgresDB, bookAddedWriter, authorRepo, bookRepoCache, bookRepo, pageRepo)
	catalogController := controller.NewCatalogController(catalogService)

	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
