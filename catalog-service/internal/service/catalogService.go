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
	GetBooksByAuthorID(authorID int) ([]model.Book, *custom.Err)
	ListBooksByAuthorName(authorName string, page, count int, sort, order string) ([]dto.ListedBooks, *custom.Err)
	ListBooksByTitle(title string, page, count int, sort, order string) ([]dto.ListedBooks, *custom.Err)
	ListBooksByAuthorNameAndTitle(authorName, title string, page, count int, sort, order string) ([]dto.ListedBooks, *custom.Err)
	ListBooksByCategory(categoryName string, page, count int, sort, order string) ([]dto.ListedBooks, *custom.Err)
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

func (s *catalogService) GetBooksByAuthorID(authorID int) ([]model.Book, *custom.Err) {
	books, err := s.bookRepository.GetBooksByAuthorID(authorID)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	return books, nil
}

func (s *catalogService) ListBooksByAuthorName(authorName string, page, count int, sort, order string) ([]dto.ListedBooks, *custom.Err) {
	rawBooks, err := s.bookRepository.ListBooksByAuthorName(authorName, page, count, sort, order)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	return s.parseListedBooksRaw(rawBooks)
}

func (s *catalogService) ListBooksByTitle(title string, page, count int, sort, order string) ([]dto.ListedBooks, *custom.Err) {
	rawBooks, err := s.bookRepository.ListBooksByTitle(title, page, count, sort, order)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	return s.parseListedBooksRaw(rawBooks)
}

func (s *catalogService) ListBooksByAuthorNameAndTitle(authorName, title string, page, count int, sort, order string) ([]dto.ListedBooks, *custom.Err) {
	rawBooks, err := s.bookRepository.ListBooksByAuthorNameAndTitle(authorName, title, page, count, sort, order)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	return s.parseListedBooksRaw(rawBooks)
}

func (s *catalogService) parseListedBooksRaw(raw []dto.ListedBooksRaw) ([]dto.ListedBooks, *custom.Err) {
	converted := make([]dto.ListedBooks, 0, len(raw))
	for _, row := range raw {
		var books []dto.ListedBook
		if err := json.Unmarshal(row.Books, &books); err != nil {
			return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
		}

		converted = append(converted, dto.ListedBooks{
			AuthorID: row.AuthorID,
			Fullname: row.Fullname,
			Books:    books,
		})
	}
	return converted, nil
}

func (s *catalogService) ListBooksByCategory(categoryName string, page, count int, sort, order string) ([]dto.ListedBooks, *custom.Err) {
	rawBooks, err := s.bookRepository.ListBooksByCategory(categoryName, page, count, sort, order)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	return s.parseListedBooksRaw(rawBooks)
}
