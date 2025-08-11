package service

import (
	"github.com/Yarik7610/library-backend/catalog-service/internal/repository"
)

type BookService interface {
}

type bookService struct {
	bookRepository repository.BookRepository
}

func NewBookService(bookRepository repository.BookRepository) BookService {
	return &bookService{bookRepository: bookRepository}
}
