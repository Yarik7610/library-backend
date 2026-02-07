package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/Yarik7610/library-backend-common/broker/kafka/event"
	"github.com/Yarik7610/library-backend/catalog-service/internal/domain"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/postgres"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/postgres/model"
	redisRepositories "github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/redis"
	postgresMapper "github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/service/mapper/postgres"
	redisMapper "github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/service/mapper/redis"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CatalogService interface {
	GetCategories(ctx context.Context) ([]string, error)
	GetNewBooks(ctx context.Context) ([]domain.Book, error)
	GetBookViewsCount(ctx context.Context, bookID uint) (int64, error)
	GetPopularBooks(ctx context.Context) ([]domain.Book, error)
	GetBooksByAuthorID(ctx context.Context, authorID uint) ([]domain.Book, error)
	GetBookPage(ctx context.Context, bookID, pageNumber uint) (*domain.Page, error)
	PreviewBook(ctx context.Context, bookID, userID uint) (*domain.Book, error)
	AddBook(ctx context.Context, bookDomain *domain.Book) error
	DeleteBook(ctx context.Context, bookID uint) error
	CreateAuthor(ctx context.Context, authorDomain *domain.Author) error
	DeleteAuthor(ctx context.Context, authorID uint) error
	ListBooksByCategory(ctx context.Context, categoryName string, page, count uint, sort, order string) ([]domain.Book, error)
	ListBooksByAuthorName(ctx context.Context, authorName string, page, count uint, sort, order string) ([]domain.Book, error)
	ListBooksByTitle(ctx context.Context, title string, page, count uint, sort, order string) ([]domain.Book, error)
	ListBooksByAuthorNameAndTitle(ctx context.Context, authorName, title string, page, count uint, sort, order string) ([]domain.Book, error)
}

type catalogService struct {
	db                       *gorm.DB
	bookAddedWriter          *kafka.Writer
	redisBookRepository      redisRepositories.BookRepository
	postgresAuthorRepository postgres.AuthorRepository
	postgresBookRepository   postgres.BookRepository
	postgresPageRepository   postgres.PageRepository
}

func NewCatalogService(
	db *gorm.DB,
	bookAddedWriter *kafka.Writer,
	redisBookRepository redisRepositories.BookRepository,
	postgresAuthorRepository postgres.AuthorRepository,
	postgresBookRepository postgres.BookRepository,
	postgresPageRepository postgres.PageRepository) CatalogService {
	return &catalogService{
		db:                       db,
		bookAddedWriter:          bookAddedWriter,
		redisBookRepository:      redisBookRepository,
		postgresAuthorRepository: postgresAuthorRepository,
		postgresBookRepository:   postgresBookRepository,
		postgresPageRepository:   postgresPageRepository,
	}
}

func (s *catalogService) GetCategories(ctx context.Context) ([]string, error) {
	categories, err := s.redisBookRepository.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	if len(categories) == 0 {
		var err error
		categories, err = s.postgresBookRepository.GetCategories(ctx)
		if err != nil {
			return nil, err
		}
		go s.redisBookRepository.SetCategories(ctx, categories)
	}

	return categories, nil
}

func (s *catalogService) GetNewBooks(ctx context.Context) ([]domain.Book, error) {
	redisNewBookWithAuthorModels, err := s.redisBookRepository.GetNew(ctx)
	if err != nil {
		return nil, err
	}

	if len(redisNewBookWithAuthorModels) > 0 {
		return redisMapper.BookWithAuthorModelsToDomains(redisNewBookWithAuthorModels), nil
	}

	postgresNewBookWithAuthorModels, err := s.postgresBookRepository.GetNew(ctx)
	if err != nil {
		return nil, err
	}

	postgresNewBookDomains := postgresMapper.BookWithAuthorModelsToDomains(postgresNewBookWithAuthorModels)
	go s.redisBookRepository.SetNew(ctx, redisMapper.BookDomainsToBookWithAuthorModels(postgresNewBookDomains))
	return postgresNewBookDomains, nil
}

func (s *catalogService) GetBookViewsCount(ctx context.Context, bookID uint) (int64, error) {
	return s.redisBookRepository.GetViewsCount(ctx, bookID)
}

func (s *catalogService) GetPopularBooks(ctx context.Context) ([]domain.Book, error) {
	bookIDs, err := s.redisBookRepository.GetPopularBookIDs(ctx)
	if err != nil {
		return nil, err
	}

	bookWithAuthorModels, err := s.postgresBookRepository.GetBooksByIDs(ctx, bookIDs)
	if err != nil {
		return nil, err
	}
	bookDomains := postgresMapper.BookWithAuthorModelsToDomains(bookWithAuthorModels)

	bookDomainsMap := make(map[uint]domain.Book)
	for _, bookDomain := range bookDomains {
		bookDomainsMap[bookDomain.ID] = bookDomain
	}

	sortedBooks := make([]domain.Book, 0)
	for _, bookIDString := range bookIDs {
		bookID, err := strconv.Atoi(bookIDString)
		if err != nil {
			return nil, err
		}
		sortedBooks = append(sortedBooks, bookDomainsMap[uint(bookID)])
	}
	return sortedBooks, nil
}

func (s *catalogService) GetBooksByAuthorID(ctx context.Context, authorID uint) ([]domain.Book, error) {
	bookWithAuthorModels, err := s.postgresBookRepository.GetBooksByAuthorID(ctx, authorID)
	if err != nil {
		return nil, err
	}
	return postgresMapper.BookWithAuthorModelsToDomains(bookWithAuthorModels), nil
}

func (s *catalogService) GetBookPage(ctx context.Context, bookID, pageNumber uint) (*domain.Page, error) {
	pageModel, err := s.postgresPageRepository.FindByBookIDAndPageNumber(ctx, bookID, pageNumber)
	if err != nil {
		return nil, err
	}
	pageDomain := postgresMapper.PageModelToDomain(pageModel)
	return &pageDomain, nil
}

func (s *catalogService) PreviewBook(ctx context.Context, bookID, userID uint) (*domain.Book, error) {
	bookWithAuthorModel, err := s.postgresBookRepository.FindByID(ctx, bookID)
	if err != nil {
		return nil, err
	}

	if userID > 0 {
		if err := s.redisBookRepository.UpdateViewsCount(ctx, bookID, userID); err != nil {
			zap.S().Warnf("Skip update book views count with ID %s because of error: %v", bookID, err)
		}
	}

	bookDomain := postgresMapper.BookWithAuthorModelToDomain(bookWithAuthorModel)
	return &bookDomain, nil
}

func (s *catalogService) AddBook(ctx context.Context, bookDomain *domain.Book) error {
	txCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var createdBookModel model.Book
	var authorModel *model.Author

	err := s.db.WithContext(txCtx).Transaction(func(tx *gorm.DB) error {
		postgresAuthorRepositoryTX := s.postgresAuthorRepository.WithinTX(tx)
		postgresPageRepositoryTX := s.postgresPageRepository.WithinTX(tx)
		postgresBookRepositoryTX := s.postgresBookRepository.WithinTX(tx)

		var err error
		authorModel, err = postgresAuthorRepositoryTX.FindByID(txCtx, bookDomain.Author.ID)
		if err != nil {
			return err
		}
		bookDomain.Author.Fullname = authorModel.Fullname

		createdBookModel = model.Book{
			AuthorID: bookDomain.Author.ID,
			Title:    bookDomain.Title,
			Year:     bookDomain.Year,
			Category: bookDomain.Category,
		}
		if err := postgresBookRepositoryTX.Create(txCtx, &createdBookModel); err != nil {
			return err
		}
		bookDomain.ID = createdBookModel.ID

		for i := range bookDomain.Pages {
			newPageModel := model.Page{
				BookID:  createdBookModel.ID,
				Number:  bookDomain.Pages[i].Number,
				Content: bookDomain.Pages[i].Content,
			}

			if err := postgresPageRepositoryTX.Create(txCtx, &newPageModel); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	bookAddedEvent, err := json.Marshal(
		event.BookAdded{
			ID:             createdBookModel.ID,
			AuthorID:       createdBookModel.AuthorID,
			AuthorFullname: authorModel.Fullname,
			Title:          createdBookModel.Title,
			Year:           createdBookModel.Year,
			Category:       createdBookModel.Category,
		})
	kafkaCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()
	if err := s.bookAddedWriter.WriteMessages(kafkaCtx, kafka.Message{Value: bookAddedEvent}); err != nil {
		zap.S().Errorf("Book added event write error: %v\n", err)
	}

	return nil
}

func (s *catalogService) DeleteBook(ctx context.Context, bookID uint) error {
	return s.postgresBookRepository.Delete(ctx, bookID)
}

func (s *catalogService) CreateAuthor(ctx context.Context, authorDomain *domain.Author) error {
	authorModel := model.Author{
		Fullname: authorDomain.Fullname,
	}

	if err := s.postgresAuthorRepository.Create(ctx, &authorModel); err != nil {
		return err
	}

	authorDomain.ID = authorModel.ID
	return nil
}

func (s *catalogService) DeleteAuthor(ctx context.Context, authorID uint) error {
	return s.postgresAuthorRepository.Delete(ctx, authorID)
}

func (s *catalogService) ListBooksByCategory(ctx context.Context, categoryName string, page, count uint, sort, order string) ([]domain.Book, error) {
	bookWithAuthorModels, err := s.postgresBookRepository.ListByCategory(ctx, categoryName, page, count, sort, order)
	if err != nil {
		return nil, err
	}
	return postgresMapper.BookWithAuthorModelsToDomains(bookWithAuthorModels), nil
}

func (s *catalogService) ListBooksByAuthorName(ctx context.Context, authorName string, page, count uint, sort, order string) ([]domain.Book, error) {
	bookWithAuthorModels, err := s.postgresBookRepository.ListByAuthorName(ctx, authorName, page, count, sort, order)
	if err != nil {
		return nil, err
	}
	return postgresMapper.BookWithAuthorModelsToDomains(bookWithAuthorModels), nil
}

func (s *catalogService) ListBooksByTitle(ctx context.Context, title string, page, count uint, sort, order string) ([]domain.Book, error) {
	bookWithAuthorModels, err := s.postgresBookRepository.ListByTitle(ctx, title, page, count, sort, order)
	if err != nil {
		return nil, err
	}
	return postgresMapper.BookWithAuthorModelsToDomains(bookWithAuthorModels), nil
}

func (s *catalogService) ListBooksByAuthorNameAndTitle(ctx context.Context, authorName, title string, page, count uint, sort, order string) ([]domain.Book, error) {
	bookWithAuthorModels, err := s.postgresBookRepository.ListByAuthorNameAndTitle(ctx, authorName, title, page, count, sort, order)
	if err != nil {
		return nil, err
	}
	return postgresMapper.BookWithAuthorModelsToDomains(bookWithAuthorModels), nil
}
