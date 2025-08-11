package service

import (
	"github.com/Yarik7610/library-backend/catalog-service/internal/repository"
)

type CatalogService interface {
}

type catalogService struct {
	bookRepository repository.BookRepository
	pageRepository repository.PageRepository
}

func NewCatalogService(bookRepository repository.BookRepository, pageRepository repository.PageRepository) CatalogService {
	return &catalogService{bookRepository: bookRepository, pageRepository: pageRepository}
}
