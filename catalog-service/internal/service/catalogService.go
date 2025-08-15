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
	GetBooksByAuthor(author string) ([]dto.Books, *custom.Err)
	GetBooksByTitle(title string) ([]dto.Books, *custom.Err)
	GetBooksByAuthorAndTitle(author, title string) ([]dto.Books, *custom.Err)
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

func (s *catalogService) GetBooksByAuthor(author string) ([]dto.Books, *custom.Err) {
	rawBooks, err := s.bookRepository.GetBooksByAuthor(author)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	return s.convertRawBooks(rawBooks)
}

func (s *catalogService) GetBooksByTitle(title string) ([]dto.Books, *custom.Err) {
	rawBooks, err := s.bookRepository.GetBooksByTitle(title)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	return s.convertRawBooks(rawBooks)
}

func (s *catalogService) GetBooksByAuthorAndTitle(author, title string) ([]dto.Books, *custom.Err) {
	rawBooks, err := s.bookRepository.GetBooksByAuthorAndTitle(author, title)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	return s.convertRawBooks(rawBooks)
}

func (s *catalogService) convertRawBooks(raw []dto.BooksRaw) ([]dto.Books, *custom.Err) {
	booksByAuthor := make([]dto.Books, 0, len(raw))
	for _, row := range raw {
		var books []dto.BookRaw
		if err := json.Unmarshal(row.Books, &books); err != nil {
			return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
		}

		booksByAuthor = append(booksByAuthor, dto.Books{
			AuthorID: row.AuthorID,
			Fullname: row.Fullname,
			Books:    books,
		})
	}
	return booksByAuthor, nil
}
