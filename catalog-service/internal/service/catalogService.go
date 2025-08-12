package service

import (
	"encoding/json"
	"net/http"

	"github.com/Yarik7610/library-backend-common/custom"
	"github.com/Yarik7610/library-backend/catalog-service/internal/dto"
	"github.com/Yarik7610/library-backend/catalog-service/internal/model"
	"github.com/Yarik7610/library-backend/catalog-service/internal/repository"
)

type CatalogService interface {
	GetCategories() ([]string, *custom.Err)
	PreviewBook(bookID uint) (*model.Book, *custom.Err)
	GetAuthorsBooks(authorName string) ([]dto.AuthorBooks, *custom.Err)
}

type catalogService struct {
	bookRepository repository.BookRepository
	pageRepository repository.PageRepository
}

func NewCatalogService(bookRepository repository.BookRepository, pageRepository repository.PageRepository) CatalogService {
	return &catalogService{bookRepository: bookRepository, pageRepository: pageRepository}
}

func (s *catalogService) GetCategories() ([]string, *custom.Err) {
	categories, err := s.bookRepository.GetCategories()
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}
	return categories, nil
}

func (s *catalogService) PreviewBook(bookID uint) (*model.Book, *custom.Err) {
	book, err := s.bookRepository.FindByID(bookID)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}
	if book == nil {
		return nil, custom.NewErr(http.StatusNotFound, "book not found")
	}
	return book, nil
}

func (s *catalogService) GetAuthorsBooks(authorName string) ([]dto.AuthorBooks, *custom.Err) {
	rawAuthorsBooks, err := s.bookRepository.GetAuthorsBooks(authorName)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	authorsBooks := make([]dto.AuthorBooks, 0, len(rawAuthorsBooks))
	for _, row := range rawAuthorsBooks {
		var books []model.Book
		if err := json.Unmarshal(row.Books, &books); err != nil {
			return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
		}

		authorsBooks = append(authorsBooks, dto.AuthorBooks{
			AuthorID: row.AuthorID,
			Fullname: row.Fullname,
			Books:    books,
		})
	}

	return authorsBooks, nil
}
