package service

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/Yarik7610/library-backend-common/broker/event"
	"github.com/Yarik7610/library-backend/catalog-service/internal/domain"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/postgres"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/postgres/model"
	redisRepositories "github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/redis"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/service/mapper"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CatalogService interface {
	GetCategories() ([]string, error)
	GetNewBooks() ([]domain.Book, error)
	GetBookViewsCount(bookID uint) (int64, error)
	GetPopularBooks() ([]domain.Book, error)
	GetBooksByAuthorID(authorID uint) ([]domain.Book, error)
	GetBookPage(bookID, pageNumber uint) (*domain.Page, error)
	PreviewBook(bookID, userID uint) (*domain.Book, error)
	AddBook(bookDomain *domain.Book) error
	DeleteBook(bookID uint) error
	CreateAuthor(authorDomain *domain.Author) error
	DeleteAuthor(authorID uint) error
	// ListBooksByCategory(categoryName string, page, count uint, sort, order string) ([]domain.ListedBooks, error)
	// ListBooksByAuthorName(authorName string, page, count uint, sort, order string) ([]domain.ListedBooks, error)
	// ListBooksByTitle(title string, page, count uint, sort, order string) ([]domain.ListedBooks, error)
	// ListBooksByAuthorNameAndTitle(authorName, title string, page, count uint, sort, order string) ([]domain.ListedBooks, error)
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

func (s *catalogService) GetCategories() ([]string, error) {
	categories, err := s.redisBookRepository.GetCategories()
	if err != nil {
		return nil, err
	}

	if len(categories) == 0 {
		var err error
		categories, err = s.postgresBookRepository.GetCategories()
		if err != nil {
			return nil, err
		}
		go s.redisBookRepository.SetCategories(categories)
	}

	return categories, nil
}

func (s *catalogService) GetNewBooks() ([]domain.Book, error) {
	newBookModels, err := s.redisBookRepository.GetNewBooks()
	if err != nil {
		return nil, err
	}

	if len(newBookModels) == 0 {
		var err error
		newBookModels, err = s.postgresBookRepository.GetNewBooks()
		if err != nil {
			return nil, err
		}
		go s.redisBookRepository.SetNewBooks(newBookModels)
	}

	return s.getBookDomains(newBookModels)
}

func (s *catalogService) GetBookViewsCount(bookID uint) (int64, error) {
	return s.redisBookRepository.GetBookViewsCount(bookID)
}

func (s *catalogService) GetPopularBooks() ([]domain.Book, error) {
	bookIDs, err := s.redisBookRepository.GetPopularBooksIDs()
	if err != nil {
		return nil, err
	}

	bookModels, err := s.postgresBookRepository.GetBooksByIDs(bookIDs)
	if err != nil {
		return nil, err
	}

	bookDomains, err := s.getBookDomains(bookModels)
	if err != nil {
		return nil, err
	}

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

func (s *catalogService) GetBooksByAuthorID(authorID uint) ([]domain.Book, error) {
	bookModels, err := s.postgresBookRepository.GetBooksByAuthorID(authorID)
	if err != nil {
		return nil, err
	}

	return s.getBookDomains(bookModels)
}

func (s *catalogService) GetBookPage(bookID, pageNumber uint) (*domain.Page, error) {
	pageModel, err := s.postgresPageRepository.GetPage(bookID, pageNumber)
	if err != nil {
		return nil, err
	}
	pageDomain := mapper.PageModelToDomain(pageModel)
	return &pageDomain, nil
}

func (s *catalogService) PreviewBook(bookID, userID uint) (*domain.Book, error) {
	bookModel, err := s.postgresBookRepository.FindByID(bookID)
	if err != nil {
		return nil, err
	}

	if userID > 0 {
		if err := s.redisBookRepository.UpdateBookViewsCount(bookID, userID); err != nil {
			zap.S().Warnf("skip update book views count with ID %s because of error: %v", bookID, err)
		}
	}

	return s.getBookDomain(bookModel)
}

func (s *catalogService) AddBook(bookDomain *domain.Book) error {
	var createdBookModel model.Book
	var authorModel *model.Author

	err := s.db.Transaction(func(tx *gorm.DB) error {
		postgresAuthorRepositoryTX := s.postgresAuthorRepository.WithinTX(tx)
		postgresPageRepositoryTX := s.postgresPageRepository.WithinTX(tx)
		postgresBookRepositoryTX := s.postgresBookRepository.WithinTX(tx)

		var err error
		authorModel, err = postgresAuthorRepositoryTX.FindByID(bookDomain.Author.ID)
		if err != nil {
			return err
		}

		createdBookModel = model.Book{
			AuthorID: bookDomain.Author.ID,
			Title:    bookDomain.Title,
			Year:     bookDomain.Year,
			Category: bookDomain.Category,
		}
		if err := postgresBookRepositoryTX.CreateBook(&createdBookModel); err != nil {
			return err
		}
		bookDomain.ID = createdBookModel.ID

		for i := range bookDomain.Pages {
			newPageModel := model.Page{
				BookID:  createdBookModel.ID,
				Number:  bookDomain.Pages[i].Number,
				Content: bookDomain.Pages[i].Content,
			}

			if err := postgresPageRepositoryTX.CreatePage(&newPageModel); err != nil {
				return err
			}
			bookDomain.Pages[i].ID = newPageModel.ID
		}

		return nil
	})
	if err != nil {
		return err
	}

	bookAddedEvent, err := json.Marshal(
		event.BookAdded{
			ID:         createdBookModel.ID,
			AuthorID:   createdBookModel.AuthorID,
			AuthorName: authorModel.Fullname,
			Title:      createdBookModel.Title,
			Year:       createdBookModel.Year,
			Category:   createdBookModel.Category,
		})
	ctx := context.Background()
	if err := s.bookAddedWriter.WriteMessages(ctx, kafka.Message{Value: bookAddedEvent}); err != nil {
		zap.S().Errorf("Book added event write error: %v\n", err)
	}

	return nil
}

func (s *catalogService) DeleteBook(bookID uint) error {
	return s.postgresBookRepository.DeleteBook(bookID)
}

func (s *catalogService) CreateAuthor(authorDomain *domain.Author) error {
	authorModel := model.Author{
		Fullname: authorDomain.Fullname,
	}

	if err := s.postgresAuthorRepository.CreateAuthor(&authorModel); err != nil {
		return err
	}

	authorDomain.ID = authorModel.ID
	return nil
}

func (s *catalogService) DeleteAuthor(authorID uint) error {
	return s.postgresAuthorRepository.DeleteAuthor(authorID)
}

// func (s *catalogService) ListBooksByCategory(categoryName string, page, count uint, sort, order string) ([]domain.ListedBooks, error) {
// 	rawBooks, err := s.postgresBookRepository.ListBooksByCategory(categoryName, page, count, sort, order)
// 	if err != nil {
// 		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
// 	}

// 	return s.parseListedBooksRaw(rawBooks)
// }

// func (s *catalogService) ListBooksByAuthorName(authorName string, page, count uint, sort, order string) ([]domain.ListedBooks, error) {
// 	rawBooks, err := s.postgresBookRepository.ListBooksByAuthorName(authorName, page, count, sort, order)
// 	if err != nil {
// 		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
// 	}

// 	return s.parseListedBooksRaw(rawBooks)
// }

// func (s *catalogService) ListBooksByTitle(title string, page, count uint, sort, order string) ([]domain.ListedBooks, error) {
// 	rawBooks, err := s.postgresBookRepository.ListBooksByTitle(title, page, count, sort, order)
// 	if err != nil {
// 		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
// 	}

// 	return s.parseListedBooksRaw(rawBooks)
// }

// func (s *catalogService) ListBooksByAuthorNameAndTitle(authorName, title string, page, count uint, sort, order string) ([]domain.ListedBooks, error) {
// 	rawBooks, err := s.postgresBookRepository.ListBooksByAuthorNameAndTitle(authorName, title, page, count, sort, order)
// 	if err != nil {
// 		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
// 	}

// 	return s.parseListedBooksRaw(rawBooks)
// }

// func (s *catalogService) parseListedBooksRaw(raw []dto.ListedBooksRaw) ([]domain.ListedBooks, error) {
// 	converted := make([]dto.ListedBooks, 0, len(raw))
// 	for _, row := range raw {
// 		var books []dto.ListedBook
// 		if err := json.Unmarshal(row.Books, &books); err != nil {
// 			return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
// 		}

// 		converted = append(converted, dto.ListedBooks{
// 			AuthorID: row.AuthorID,
// 			Fullname: row.Fullname,
// 			Books:    books,
// 		})
// 	}
// 	return converted, nil
// }

func (s *catalogService) getBookDomains(bookModels []model.Book) ([]domain.Book, error) {
	bookDomains := make([]domain.Book, len(bookModels))

	for i, bookModel := range bookModels {
		bookModel, err := s.getBookDomain(&bookModel)
		if err != nil {
			return nil, err
		}
		bookDomains[i] = *bookModel
	}
	return bookDomains, nil
}

func (s *catalogService) getBookDomain(bookModel *model.Book) (*domain.Book, error) {
	authorModel, err := s.postgresAuthorRepository.FindByID(bookModel.AuthorID)
	if err != nil {
		return nil, err
	}

	bookDomain := mapper.BookModelToDomain(bookModel)
	bookDomain.Author = mapper.AuthorModelToDomain(authorModel)
	return &bookDomain, nil
}
