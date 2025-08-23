package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Yarik7610/library-backend-common/custom"
	"github.com/Yarik7610/library-backend/catalog-service/internal/dto"
	"github.com/Yarik7610/library-backend/catalog-service/internal/model"
	"github.com/Yarik7610/library-backend/catalog-service/internal/repository"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type CatalogService interface {
	GetCategories() ([]string, *custom.Err)
	GetNewBooks() ([]model.Book, *custom.Err)
	ListBooksByCategory(categoryName string, page, count uint, sort, order string) ([]dto.ListedBooks, *custom.Err)
	GetBooksByAuthorID(authorID uint) ([]model.Book, *custom.Err)
	ListBooksByAuthorName(authorName string, page, count uint, sort, order string) ([]dto.ListedBooks, *custom.Err)
	ListBooksByTitle(title string, page, count uint, sort, order string) ([]dto.ListedBooks, *custom.Err)
	ListBooksByAuthorNameAndTitle(authorName, title string, page, count uint, sort, order string) ([]dto.ListedBooks, *custom.Err)
	PreviewBook(bookID uint) (*model.Book, *custom.Err)
	GetBookPage(bookID, pageNumber uint) (*model.Page, *custom.Err)
	DeleteBook(bookID uint) *custom.Err
	AddBook(book *dto.AddBook) (*model.Book, *custom.Err)
	DeleteAuthor(authorID uint) *custom.Err
	CreateAuthor(fullname string) (*model.Author, *custom.Err)
}

type catalogService struct {
	db                  *gorm.DB
	authorRepository    repository.AuthorRepository
	bookRepositoryCache repository.BookRepositoryCache
	bookRepository      repository.BookRepository
	pageRepository      repository.PageRepository
}

func NewCatalogService(db *gorm.DB,
	authorRepository repository.AuthorRepository,
	bookRepositoryCache repository.BookRepositoryCache,
	bookRepository repository.BookRepository,
	pageRepository repository.PageRepository) CatalogService {
	return &catalogService{
		db:                  db,
		authorRepository:    authorRepository,
		bookRepositoryCache: bookRepositoryCache,
		bookRepository:      bookRepository,
		pageRepository:      pageRepository,
	}
}

func (s *catalogService) GetCategories() ([]string, *custom.Err) {
	categories, err := s.bookRepositoryCache.GetCategories()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	if len(categories) == 0 {
		var err error
		categories, err = s.bookRepository.GetCategories()
		if err != nil {
			return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
		}
		go s.bookRepositoryCache.SetCategories(categories)
	}

	return categories, nil
}

func (s *catalogService) GetNewBooks() ([]model.Book, *custom.Err) {
	newBooks, err := s.bookRepositoryCache.GetNewBooks()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	if len(newBooks) == 0 {
		var err error
		newBooks, err = s.bookRepository.GetNewBooks()
		if err != nil {
			return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
		}
		go s.bookRepositoryCache.SetNewBooks(newBooks)
	}

	return newBooks, nil
}

func (s *catalogService) PreviewBook(bookID uint) (*model.Book, *custom.Err) {
	book, err := s.bookRepository.FindByID(bookID)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}
	if book == nil {
		return nil, custom.NewErr(http.StatusNotFound, fmt.Sprintf("book with ID %d not found", bookID))
	}
	return book, nil
}

func (s *catalogService) GetBooksByAuthorID(authorID uint) ([]model.Book, *custom.Err) {
	books, err := s.bookRepository.GetBooksByAuthorID(authorID)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	return books, nil
}

func (s *catalogService) ListBooksByAuthorName(authorName string, page, count uint, sort, order string) ([]dto.ListedBooks, *custom.Err) {
	rawBooks, err := s.bookRepository.ListBooksByAuthorName(authorName, page, count, sort, order)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	return s.parseListedBooksRaw(rawBooks)
}

func (s *catalogService) ListBooksByTitle(title string, page, count uint, sort, order string) ([]dto.ListedBooks, *custom.Err) {
	rawBooks, err := s.bookRepository.ListBooksByTitle(title, page, count, sort, order)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	return s.parseListedBooksRaw(rawBooks)
}

func (s *catalogService) ListBooksByAuthorNameAndTitle(authorName, title string, page, count uint, sort, order string) ([]dto.ListedBooks, *custom.Err) {
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

func (s *catalogService) ListBooksByCategory(categoryName string, page, count uint, sort, order string) ([]dto.ListedBooks, *custom.Err) {
	rawBooks, err := s.bookRepository.ListBooksByCategory(categoryName, page, count, sort, order)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	return s.parseListedBooksRaw(rawBooks)
}

func (s *catalogService) GetBookPage(bookID, pageNumber uint) (*model.Page, *custom.Err) {
	page, err := s.pageRepository.GetPage(bookID, pageNumber)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}
	if page == nil {
		return nil, custom.NewErr(http.StatusNotFound, fmt.Sprintf("book with ID %d and page number %d not found", bookID, pageNumber))
	}
	return page, nil
}

func (s *catalogService) DeleteBook(bookID uint) *custom.Err {
	err := s.bookRepository.DeleteBook(bookID)
	if err != nil {
		return custom.NewErr(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (s *catalogService) AddBook(book *dto.AddBook) (*model.Book, *custom.Err) {
	var created model.Book

	err := s.db.Transaction(func(tx *gorm.DB) error {
		authorRepositoryTX := s.authorRepository.WithinTX(tx)
		pageRepositoryTX := s.pageRepository.WithinTX(tx)
		bookRepositoryTX := s.bookRepository.WithinTX(tx)

		author, err := authorRepositoryTX.FindByID(book.AuthorID)
		if err != nil {
			return err
		}
		if author == nil {
			return fmt.Errorf("author with ID %d doesn't exist, create author first", book.AuthorID)
		}

		foundBook, err := bookRepositoryTX.FindByTitleAndAuthorID(book.Title, book.AuthorID)
		if err != nil {
			return nil
		}
		if foundBook != nil {
			return fmt.Errorf("book with author ID %d and title %s already exists", foundBook.AuthorID, foundBook.Title)
		}

		created = model.Book{
			AuthorID: book.AuthorID,
			Title:    book.Title,
			Year:     book.Year,
			Category: book.Category,
		}
		err = bookRepositoryTX.CreateBook(&created)
		if err != nil {
			return err
		}

		for _, page := range book.Pages {
			newPage := model.Page{
				BookID:  created.ID,
				Number:  page.Number,
				Content: page.Content,
			}
			err := pageRepositoryTX.CreatePage(&newPage)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	return &created, nil
}

func (s *catalogService) DeleteAuthor(authorID uint) *custom.Err {
	err := s.authorRepository.DeleteAuthor(authorID)
	if err != nil {
		return custom.NewErr(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (s *catalogService) CreateAuthor(fullname string) (*model.Author, *custom.Err) {
	author := model.Author{
		Fullname: fullname,
	}
	err := s.authorRepository.CreateAuthor(&author)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}
	return &author, nil
}
