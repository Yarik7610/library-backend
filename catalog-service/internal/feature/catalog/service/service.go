package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Yarik7610/library-backend-common/broker/event"
	"github.com/Yarik7610/library-backend-common/custom"
	"github.com/Yarik7610/library-backend/catalog-service/internal/domain"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/postgres"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/postgres/model"
	redisRepositories "github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/redis"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/transport/http/dto"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CatalogService interface {
	GetCategories() ([]string, *custom.Err)
	GetNewBooks() ([]domain.Book, *custom.Err)
	GetBookViewsCount(bookID uint) (int64, *custom.Err)
	GetPopularBooks() ([]domain.Book, *custom.Err)
	ListBooksByCategory(categoryName string, page, count uint, sort, order string) ([]domain.ListedBooks, *custom.Err)
	GetBooksByAuthorID(authorID uint) ([]domain.Book, *custom.Err)
	ListBooksByAuthorName(authorName string, page, count uint, sort, order string) ([]domain.ListedBooks, *custom.Err)
	ListBooksByTitle(title string, page, count uint, sort, order string) ([]domain.ListedBooks, *custom.Err)
	ListBooksByAuthorNameAndTitle(authorName, title string, page, count uint, sort, order string) ([]domain.ListedBooks, *custom.Err)
	PreviewBook(bookID, userID uint) (*domain.Book, *custom.Err)
	GetBookPage(bookID, pageNumber uint) (*domain.Page, *custom.Err)
	DeleteBook(bookID uint) *custom.Err
	AddBook(book *domain.Book) *custom.Err
	DeleteAuthor(authorID uint) *custom.Err
	CreateAuthor(fullname string) (*domain.Author, *custom.Err)
}

type catalogService struct {
	db                       *gorm.DB
	bookAddedWriter          *kafka.Writer
	postgresAuthorRepository postgres.AuthorRepository
	redisBookRepository      redisRepositories.BookRepository
	postgresBookRepository   postgres.BookRepository
	postgresPageRepository   postgres.PageRepository
}

func NewCatalogService(
	db *gorm.DB,
	bookAddedWriter *kafka.Writer,
	postgresAuthorRepository postgres.AuthorRepository,
	redisBookRepository redisRepositories.BookRepository,
	postgresBookRepository postgres.BookRepository,
	postgresPageRepository postgres.PageRepository) CatalogService {
	return &catalogService{
		db:                       db,
		bookAddedWriter:          bookAddedWriter,
		postgresAuthorRepository: postgresAuthorRepository,
		redisBookRepository:      redisBookRepository,
		postgresBookRepository:   postgresBookRepository,
		postgresPageRepository:   postgresPageRepository,
	}
}

func (s *catalogService) GetCategories() ([]string, *custom.Err) {
	categories, err := s.redisBookRepository.GetCategories()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	if len(categories) == 0 {
		var err error
		categories, err = s.postgresBookRepository.GetCategories()
		if err != nil {
			return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
		}
		go s.redisBookRepository.SetCategories(categories)
	}

	return categories, nil
}

func (s *catalogService) GetNewBooks() ([]domain.Book, *custom.Err) {
	newBookModels, err := s.redisBookRepository.GetNewBooks()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	if len(newBookModels) == 0 {
		var err error
		newBookModels, err = s.postgresBookRepository.GetNewBooks()
		if err != nil {
			return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
		}
		go s.redisBookRepository.SetNewBooks(newBookModels)
	}

	return newBookModels, nil
}

func (s *catalogService) GetBookViewsCount(bookID uint) (int64, *custom.Err) {
	viewsCount, err := s.redisBookRepository.GetBookViewsCount(bookID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, custom.NewErr(http.StatusInternalServerError, err.Error())
	}
	return viewsCount, nil
}

func (s *catalogService) GetPopularBooks() ([]model.Book, *custom.Err) {
	booksIDs, err := s.redisBookRepository.GetPopularBooksIDs()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	books, err := s.postgresBookRepository.GetBooksByIDs(booksIDs)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	booksMap := make(map[uint]model.Book)
	for _, b := range books {
		booksMap[b.ID] = b
	}

	sortedBooks := make([]model.Book, 0)
	for _, bookIDString := range booksIDs {
		bookID, err := strconv.Atoi(bookIDString)
		if err != nil {
			return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
		}
		sortedBooks = append(sortedBooks, booksMap[uint(bookID)])
	}

	return sortedBooks, nil
}

func (s *catalogService) PreviewBook(bookID, userID uint) (*model.Book, *custom.Err) {
	book, err := s.postgresBookRepository.FindByID(bookID)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}
	if book == nil {
		return nil, custom.NewErr(http.StatusNotFound, fmt.Sprintf("book with ID %d not found", bookID))
	}

	if userID > 0 {
		if err := s.redisBookRepository.UpdateBookViewsCount(bookID, userID); err != nil {
			zap.S().Warnf("skip update book views count with ID %s because of error: %v", bookID, err)
		}
	}

	return book, nil
}

func (s *catalogService) GetBooksByAuthorID(authorID uint) ([]domain.Book, *custom.Err) {
	bookModels, err := s.postgresBookRepository.GetBooksByAuthorID(authorID)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	bookDomains := make([]domain.Book, len(bookModels))
	for i, bookModel := range bookModels {
		bookDomains[i] = bookModel.ToDomain()
	}
	return bookDomains, nil
}

func (s *catalogService) ListBooksByAuthorName(authorName string, page, count uint, sort, order string) ([]dto.ListedBooks, *custom.Err) {
	rawBooks, err := s.postgresBookRepository.ListBooksByAuthorName(authorName, page, count, sort, order)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	return s.parseListedBooksRaw(rawBooks)
}

func (s *catalogService) ListBooksByTitle(title string, page, count uint, sort, order string) ([]dto.ListedBooks, *custom.Err) {
	rawBooks, err := s.postgresBookRepository.ListBooksByTitle(title, page, count, sort, order)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	return s.parseListedBooksRaw(rawBooks)
}

func (s *catalogService) ListBooksByAuthorNameAndTitle(authorName, title string, page, count uint, sort, order string) ([]dto.ListedBooks, *custom.Err) {
	rawBooks, err := s.postgresBookRepository.ListBooksByAuthorNameAndTitle(authorName, title, page, count, sort, order)
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
	rawBooks, err := s.postgresBookRepository.ListBooksByCategory(categoryName, page, count, sort, order)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	return s.parseListedBooksRaw(rawBooks)
}

func (s *catalogService) GetBookPage(bookID, pageNumber uint) (*domain.Page, *custom.Err) {
	pageModel, err := s.postgresPageRepository.GetPage(bookID, pageNumber)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}
	if pageModel == nil {
		return nil, custom.NewErr(http.StatusNotFound, fmt.Sprintf("book with ID %d and page number %d not found", bookID, pageNumber))
	}
	return pageModel.ToDomain(), nil
}

func (s *catalogService) DeleteBook(bookID uint) *custom.Err {
	err := s.postgresBookRepository.DeleteBook(bookID)
	if err != nil {
		return custom.NewErr(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (s *catalogService) AddBook(bookDomain *domain.Book) *custom.Err {
	var createdBookModel model.Book
	var author *model.Author

	err := s.db.Transaction(func(tx *gorm.DB) error {
		postgresAuthorRepositoryTX := s.postgresAuthorRepository.WithinTX(tx)
		postgresPageRepositoryTX := s.postgresPageRepository.WithinTX(tx)
		postgresBookRepositoryTX := s.postgresBookRepository.WithinTX(tx)

		var err error
		author, err = postgresAuthorRepositoryTX.FindByID(bookDomain.Author.ID)
		if err != nil {
			return err
		}
		if author == nil {
			return fmt.Errorf("author with ID %d doesn't exist, create author first", bookDomain.Author.ID)
		}

		foundBook, err := postgresBookRepositoryTX.FindByTitleAndAuthorID(bookDomain.Title, bookDomain.Author.ID)
		if err != nil {
			return nil
		}
		if foundBook != nil {
			return fmt.Errorf("book with author ID %d and title %s already exists", foundBook.AuthorID, foundBook.Title)
		}

		createdBookModel = model.Book{
			AuthorID: bookDomain.Author.ID,
			Title:    bookDomain.Title,
			Year:     bookDomain.Year,
			Category: bookDomain.Category,
		}
		err = postgresBookRepositoryTX.CreateBook(&createdBookModel)
		if err != nil {
			return err
		}

		for _, page := range bookDomain.Pages {
			newPage := model.Page{
				BookID:  createdBookModel.ID,
				Number:  page.Number,
				Content: page.Content,
			}
			err := postgresPageRepositoryTX.CreatePage(&newPage)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	bookAddedEvent, err := json.Marshal(
		event.BookAdded{
			ID:         createdBookModel.ID,
			AuthorID:   createdBookModel.AuthorID,
			AuthorName: author.Fullname,
			Title:      createdBookModel.Title,
			Year:       createdBookModel.Year,
			Category:   createdBookModel.Category,
		})
	ctx := context.Background()
	if err := s.bookAddedWriter.WriteMessages(ctx, kafka.Message{Value: bookAddedEvent}); err != nil {
		zap.S().Errorf("Book added event write error: %v\n", err)
	}

	// ID        uint   `gorm:"primarykey"`
	//   AuthorID  uint   `gorm:"uniqueIndex:author_id_title_index"`
	//   Title     string `gorm:"uniqueIndex:author_id_title_index"`
	//   Year      int
	//   Category  string
	//   CreatedAt time.Time
	//   Pages     []Page `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// bookDomain.

	//COPY every page and author details to domain, probably in TX

	return nil
}

func (s *catalogService) DeleteAuthor(authorID uint) *custom.Err {
	err := s.postgresAuthorRepository.DeleteAuthor(authorID)
	if err != nil {
		return custom.NewErr(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (s *catalogService) CreateAuthor(fullname string) (*domain.Author, *custom.Err) {
	authorModel := model.Author{
		Fullname: fullname,
	}
	err := s.postgresAuthorRepository.CreateAuthor(&authorModel)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	return &domain.Author{
		ID:       authorModel.ID,
		Fullname: authorModel.Fullname,
	}, nil
}
